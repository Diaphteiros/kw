package bookmark

import (
	"github.com/spf13/cobra"
)

var ClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"cl", "c"},
	Args:    cobra.NoArgs,
	Short:   "Delete all bookmark entries",
	Long: `Delete all bookmark entries.

This is an alias for 'forget --all'.

Note that the bookmarks are shared between all terminal sessions.`,
	Run: func(cmd *cobra.Command, args []string) {
		all = true
		ForgetCmd.Run(cmd, args)
		cmd.Println("Bookmarks cleared.")
	},
}
