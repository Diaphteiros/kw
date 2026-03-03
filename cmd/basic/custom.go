package basic

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
	"github.com/Diaphteiros/kw/pkg/config"

	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

const CustomCmdPluginName = "<internal_custom>"

// CustomCmd represents the custom command
var CustomCmd = &cobra.Command{
	Use:     "custom <path>",
	Aliases: []string{"cus", "c"},
	Args:    cobra.ExactArgs(1),
	GroupID: cmdgroups.Basic,
	Short:   "Switch to the kubeconfig at the specified path",
	Long: `Switch to the kubeconfig at the specified path.

This creates a symlink from $KUBECONFIG pointing to the given file.`,
	Run: func(cmd *cobra.Command, args []string) {
		con := config.Runtime.Context()
		con.SetPluginName(CustomCmdPluginName)

		path := args[0]
		path, err := filepath.Abs(path)
		if err != nil {
			libutils.Fatal(1, "error getting absolute path: %w\n", err)
		}
		if err := con.WriteKubeconfigSymlink(path, "Switched to custom kubeconfig from '%s'.", path); err != nil {
			libutils.Fatal(1, "error creating kubeconfig symlink: %w\n", err)
		}
		// shorten path for id if a max id length is set
		pathForId := path
		customCfg := config.Runtime.Config().Builtin.GetBuiltinCustomConfig()
		if customCfg != nil && customCfg.MaxIdLength > 0 && len(pathForId) > customCfg.MaxIdLength {
			pathForId = fmt.Sprintf("…%s", path[len(path)-customCfg.MaxIdLength:])
		}
		if err := con.WriteId("custom:%s", pathForId); err != nil {
			libutils.Fatal(1, "error writing id: %w\n", err)
		}
		if err := con.WritePluginState(&customPluginState{Source: path}); err != nil {
			libutils.Fatal(1, "error writing plugin state: %w\n", err)
		}
	},
}

type customPluginState struct {
	Source string `json:"source"`
}
