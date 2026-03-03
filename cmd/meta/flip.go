package meta

import (
	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/cmd/meta/history"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
)

var FlipCmd = &cobra.Command{
	Use:     "flip",
	Aliases: []string{"f"},
	Args:    cobra.NoArgs,
	GroupID: cmdgroups.Meta,
	Short:   "Flip the current configuration with the previously used one",
	Long: `Flip the current configuration with the previously used one.

This is basically an alias for 'kw history load 1'.
It requires a history depth of at least 1 (meaning the history must not be deactivated).`,
	Run: func(cmd *cobra.Command, args []string) {
		history.HistoryLoadCmd.Run(cmd, []string{"1"})
	},
}
