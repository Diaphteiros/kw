package storage

import (
	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/cmd/meta/bookmark"
)

var (
	all bool
)

var ForgetCmd = &cobra.Command{
	Use:     "forget [<name> ...]",
	Aliases: []string{"f"},
	Args:    cobra.ArbitraryArgs,
	Short:   "[DEPRECATED] Delete a storage entry",
	Long: `Delete a storage entry.
	
Either one or more storage keys or the '--all' flag must be specified.
If storage keys are specified, the corresponding storage entries will be deleted.
Missing entries are ignored.

If the '--all' flag is set, all storage entries will be deleted.

Note that the storage is shared between all terminal sessions.`,
	Deprecated: "use 'bookmark forget' instead.",
	Run:        bookmark.ForgetCmd.Run,
}

func init() {
	ForgetCmd.Flags().BoolVar(&all, "all", false, "Delete all storage entries")
}
