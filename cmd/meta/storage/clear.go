package storage

import (
	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/cmd/meta/bookmark"
)

var ClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"cl", "c"},
	Args:    cobra.NoArgs,
	Short:   "[DEPRECATED] Delete all storage entries",
	Long: `Delete all storage entries.

This is an alias for 'forget --all'.

Note that the storage is shared between all terminal sessions.`,
	Deprecated: "use 'bookmark clear' instead.",
	Run:        bookmark.ClearCmd.Run,
}
