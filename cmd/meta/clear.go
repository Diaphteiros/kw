package meta

import (
	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/cmdgroups"
	"github.com/Diaphteiros/kw/pkg/config"

	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var ClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"cl"},
	Args:    cobra.NoArgs,
	GroupID: cmdgroups.Meta,
	Short:   "Clear the current session's state",
	Long: `Clear the current session's state.

This removes the entire session directory, wiping the active kubeconfig,
state files, history, and any other session-specific data.`,
	Run: func(cmd *cobra.Command, args []string) {
		sdir := config.Runtime.SessionDir()
		if err := fs.FS.RemoveAll(sdir); err != nil {
			libutils.Fatal(1, "error removing session directory '%s': %s\n", sdir, err.Error())
		}
		cmd.Println("Session state cleared.")
	},
}
