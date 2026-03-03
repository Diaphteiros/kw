package meta

import (
	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/cmd/meta/storage"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
)

var StorageCmd = &cobra.Command{
	Use:     "storage",
	Aliases: []string{"store", "s"},
	GroupID: cmdgroups.Meta,
	Short:   "[DEPRECATED] Interact with the kubeconfig storage",
	Long: `Interact with the kubeconfig storage.

The storage can be used to store and load kubeconfigs.
See the different subcommands for more information.

Note that the storage is shared between all terminal sessions.`,
	Deprecated: "use 'bookmark' instead.",
	Run:        BookmarkCmd.Run,
}

func init() {
	StorageCmd.AddCommand(storage.StoreCmd)
	StorageCmd.AddCommand(storage.LoadCmd)
	StorageCmd.AddCommand(storage.ClearCmd)
	StorageCmd.AddCommand(storage.ViewCmd)
	StorageCmd.AddCommand(storage.ForgetCmd)
}
