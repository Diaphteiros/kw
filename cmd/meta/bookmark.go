package meta

import (
	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/cmd/meta/bookmark"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
)

var BookmarkCmd = &cobra.Command{
	Use:     "bookmark",
	Aliases: []string{"book", "bm", "b"},
	GroupID: cmdgroups.Meta,
	Short:   "Interact with the kubeconfig bookmarks",
	Long: `Interact with the kubeconfig bookmarks.

The bookmarks can be used to store and load kubeconfigs.
See the different subcommands for more information.

Note that the bookmarks are shared between all terminal sessions.`,
}

func init() {
	BookmarkCmd.AddCommand(bookmark.SaveCmd)
	BookmarkCmd.AddCommand(bookmark.LoadCmd)
	BookmarkCmd.AddCommand(bookmark.ClearCmd)
	BookmarkCmd.AddCommand(bookmark.ViewCmd)
	BookmarkCmd.AddCommand(bookmark.ForgetCmd)
}
