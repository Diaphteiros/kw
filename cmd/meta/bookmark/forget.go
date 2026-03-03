package bookmark

import (
	"path/filepath"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/pkg/storage"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var (
	all bool
)

var ForgetCmd = &cobra.Command{
	Use:     "forget [<name> ...]",
	Aliases: []string{"f"},
	Args:    cobra.ArbitraryArgs,
	Short:   "Delete a bookmark entry",
	Long: `Delete a bookmark entry.
	
Either one or more bookmark keys or the '--all' flag must be specified.
If bookmark keys are specified, the corresponding bookmark entries will be deleted.
Missing entries are ignored.

If the '--all' flag is set, all bookmark entries will be deleted.

Note that the bookmarks are shared between all terminal sessions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if (len(args) > 0) == all {
			libutils.Fatal(1, "either one or more bookmark keys or the '--all' flag must be specified")
		}

		sRoot := storage.GetStorageRootPath()
		var names []string
		if all {
			sDirs, err := vfs.ReadDir(fs.FS, sRoot)
			if err != nil && !vfs.IsNotExist(err) {
				libutils.Fatal(1, "error reading bookmark directory: %w\n", err)
			}
			names = make([]string, 0, len(sDirs))
			for _, sDir := range sDirs {
				if !sDir.IsDir() {
					continue
				}
				names = append(names, sDir.Name())
			}
		} else {
			names = args
		}
		for _, name := range names {
			debug.Debug("Deleting bookmark entry '%s'", name)
			if err := fs.FS.RemoveAll(filepath.Join(sRoot, name)); err != nil {
				libutils.Fatal(1, "error deleting bookmark entry '%s': %w", name, err)
			}
			cmd.Printf("Bookmark entry '%s' deleted.\n", name)
		}
	},
}

func init() {
	ForgetCmd.Flags().BoolVar(&all, "all", false, "Delete all bookmark entries")
}
