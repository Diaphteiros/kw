package history

import (
	"path/filepath"
	"strings"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/storage"

	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var HistoryClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"cl", "c"},
	Args:    cobra.NoArgs,
	Short:   "Clear the history",
	Long: `Clear the history.
	
This removes the entire history.`,
	Run: func(cmd *cobra.Command, args []string) {
		hRoot := storage.GetHistoryRootPath(Global)

		// in local mode, we have to copy entries that are still referenced by global history entries
		if !Global {
			// render global history
			debug.Debug("Checking for global history references")
			ghRoot := storage.GetHistoryRootPath(true)
			hDirs, err := vfs.ReadDir(fs.FS, ghRoot)
			if err != nil && !vfs.IsNotExist(err) {
				libutils.Fatal(1, "error reading global history directory '%s': %s\n", ghRoot, err.Error())
			}
			// check for symlinks to local history entries
			for _, hDir := range hDirs {
				hDirPath := filepath.Join(ghRoot, hDir.Name())
				if hDir.IsDir() {
					debug.Debug("Skipping global history file '%s' because it is a directory", hDirPath)
					continue
				}
				realDir, err := fs.FS.Readlink(hDirPath)
				if err != nil {
					debug.Debug("Ignoring global history entry '%s' because it is not a symlink", hDirPath)
					continue
				}
				if strings.HasPrefix(realDir, hRoot) {
					debug.Debug("Found global history reference to local history entry '%s'", realDir)
					// remove the global history entry
					if err := fs.FS.Remove(hDirPath); err != nil {
						libutils.Fatal(1, "error removing global history entry '%s': %s\n", hDirPath, err.Error())
					}
					// move local directory to global history directory
					if err := fs.FS.Rename(realDir, hDirPath); err != nil {
						libutils.Fatal(1, "error moving local history entry '%s' to global history directory: %s\n", realDir, err.Error())
					}
				}
			}
		}

		debug.Debug("Removing history directory '%s'", hRoot)
		if err := fs.FS.RemoveAll(hRoot); err != nil {
			libutils.Fatal(1, "error removing history directory '%s': %s\n", hRoot, err.Error())
		}
		cmd.Println("History cleared.")
	},
}
