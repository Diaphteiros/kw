package config

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/pkg/config"
)

var dir bool

var PathCmd = &cobra.Command{
	Use:     "path",
	Aliases: []string{"p"},
	Args:    cobra.NoArgs,
	Short:   "View the path of the kubeswitcher configuration file",
	Long:    `View the path of the kubeswitcher configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := config.Runtime.ConfigDirectory()
		if !dir {
			path = filepath.Join(path, config.ConfigFileName)
		}
		cmd.Println(path)
	},
}

func init() {
	PathCmd.Flags().BoolVarP(&dir, "dir", "d", false, "Print the path to the configuration directory instead of the configuration file")
}
