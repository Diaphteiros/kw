package meta

import (
	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/cmd/meta/history"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
)

var HistoryCmd = &cobra.Command{
	Use:     "history <command>",
	Aliases: []string{"hist", "h"},
	GroupID: cmdgroups.Meta,
	Short:   "Interact with the history",
	Long: `Interact with the history.

The subcommands allow to view the history and load a specific entry again.`,
}

func init() {
	HistoryCmd.AddCommand(history.HistoryViewCmd)
	HistoryCmd.AddCommand(history.HistoryClearCmd)
	HistoryCmd.AddCommand(history.HistoryLoadCmd)

	HistoryCmd.PersistentFlags().BoolVarP(&history.Global, "global", "g", false, "Use the global history instead of the session-specific one")
}
