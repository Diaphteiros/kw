package basic

import (
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/cmdgroups"
	"github.com/Diaphteiros/kw/pkg/config"

	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

const TemporaryCmdPluginName = "<internal_temporary>"

// TemporaryCmd represents the temporary command
var TemporaryCmd = &cobra.Command{
	Use:     "temporary",
	Aliases: []string{"temp", "tmp", "t"},
	Args:    cobra.NoArgs,
	GroupID: cmdgroups.Basic,
	Short:   "Switch to the kubeconfig currently contained in the clipboard",
	Long:    `Switch to the kubeconfig currently contained in the clipboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		con := config.Runtime.Context()
		con.SetPluginName(TemporaryCmdPluginName)

		kcfg_string, err := clipboard.ReadAll()
		if err != nil {
			libutils.Fatal(1, "error reading from clipboard: %w\n", err)
		}
		kcfg_data := []byte(kcfg_string)
		if err := con.WriteKubeconfig(kcfg_data, "Switched to kubeconfig from clipboard."); err != nil {
			libutils.Fatal(1, "error writing kubeconfig: %w\n", err)
		}
		kcfg, err := libutils.ParseKubeconfig(kcfg_data)
		if err != nil {
			libutils.Fatal(1, "error parsing kubeconfig: %w\n", err)
		}
		host, err := libutils.GetCurrentApiserverHost(kcfg)
		if err != nil {
			libutils.Fatal(1, "error getting current apiserver host from kubeconfig: %w\n", err)
		}
		if err := con.WriteId("temp:%s", host); err != nil {
			libutils.Fatal(1, "error writing id: %w\n", err)
		}
	},
}
