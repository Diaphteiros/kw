package meta

import (
	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/cmd/meta/history"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
)

var RepeatCmd = &cobra.Command{
	Use:     "repeat",
	Aliases: []string{"r"},
	Args:    cobra.NoArgs,
	GroupID: cmdgroups.Meta,
	Short:   "Switch to the last used configuration (cross-session)",
	Long: `Switch to the last used configuration (cross-session).

This is basically an alias for 'kw history load 0 --global'.
It requires a history depth of at least 1 (meaning the history must not be deactivated).`,
	Run: func(cmd *cobra.Command, args []string) {
		history.Global = true
		history.HistoryLoadCmd.Run(cmd, []string{"0"})
	},
}
