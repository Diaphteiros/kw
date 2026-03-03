package config

import (
	"github.com/spf13/cobra"

	configsubcommands "github.com/Diaphteiros/kw/cmd/config/config"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
)

var ConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Args:    cobra.NoArgs,
	GroupID: cmdgroups.Config,
	Short:   "Interact with the kubeswitcher configuration",
	Long: `Interact with the kubeswitcher configuration.

The default configuration directory is $HOME/.kubeswitcher_config/.
This can be changed by setting the environment variable KW_CONFIG_REPO.

Within the configuration directory, a 'config.yaml' file is expected to contain the kubeswitcher configuration.
Directory and file (with default values) will be created when any kubeswitcher command is run, if they do not exist.`,
}

func init() {
	ConfigCmd.AddCommand(configsubcommands.PathCmd)
}
