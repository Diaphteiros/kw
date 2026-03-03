package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	StoreSubPath           = "store"
	HistroySubPath         = "history"
	CurrentHistoryFileName = "current"
	HistoryTmpKey          = "tmp"
)

var (
	allFileMappings = []*FilenameMapping{Fnmi(config.KubeconfigFileName).WithPossibleSymlink(), Fnmi(config.GenericStateFileName), Fnmi(config.PluginStateFileName), Fnm(config.NotificationMessageBackupFileName, config.NotificationMessageFileName), Fnmi(config.NotificationMessageFileName), Fnmi(config.IdFileName)} // store backup files as the actual files, unless the actual files exist
)

// GetStoragePath returns the path to the storage directory for the given id.
func GetStoragePath(id string) string {
	return filepath.Join(GetStorageRootPath(), id)
}

// GetStorageRootPath returns the path to the root directory for the storage.
func GetStorageRootPath() string {
	return filepath.Join(config.Runtime.ConfigDirectory(), StoreSubPath)
}

// GetHistoryPath returns the path to the history directory for the given id.
func GetHistoryPath(id string, global bool) string {
	return filepath.Join(GetHistoryRootPath(global), id)
}

// GetHistoryRootPath returns the path to the root directory for the history.
func GetHistoryRootPath(global bool) string {
	if global {
		return filepath.Join(config.Runtime.ConfigDirectory(), HistroySubPath)
	}
	return filepath.Join(config.Runtime.SessionDir(), HistroySubPath)
}

// GetCurrentHistoryIndexFilePath returns the path to the file containing the current history index.
func GetCurrentHistoryIndexFilePath(global bool) string {
	return filepath.Join(GetHistoryRootPath(global), CurrentHistoryFileName)
}

// GetCurrentHistoryIndex returns the current history index.
// Returns -1 if the file doesn't exist and an error if the file can't be read or its contents don't parse into an integer.
func GetCurrentHistoryIndex(global bool) (int, error) {
	rawCurrentHistoryIndex, err := vfs.ReadFile(fs.FS, GetCurrentHistoryIndexFilePath(global))
	if err != nil {
		if vfs.IsNotExist(err) {
			return -1, nil
		}
		return -1, fmt.Errorf("error reading current history index (global=%t): %w", global, err)
	}
	currentHistoryIndex, err := strconv.Atoi(string(rawCurrentHistoryIndex))
	if err != nil {
		return -1, fmt.Errorf("error parsing current history index (global=%t): %w", global, err)
	}
	return currentHistoryIndex, nil
}

// ConvertHistoryIndexAndDirName converts between the current history index and a history directory name (which is an index too).
// This is necessary because the history directories are used like a cyclic buffer.
// The function is symmetrical and can be used to convert in both directions.
// If the returned index is negative or equal to/greater than the history depth, it means that either the directory is not part of the history or no corresponding directory was found.
// This can happen as a result of a changed history depth.
func ConvertHistoryIndexAndDirName(currentHistoryIndex, indexDirName int) int {
	actual := currentHistoryIndex - indexDirName
	if actual < 0 {
		actual += config.Runtime.Config().Kubeswitcher.HistoryDepth
	}
	return actual
}

// StoreOrLoadFiles takes a source and destination path and a mapping of filenames
// For each f1:f2 mapping, it copies 'src/f1' to 'dst/f2', if 'src/f1' exists.
// Otherwise, it removes 'dst/f2', if it exists.
// If multiple mappings point to the same destination file, latter mappings will overwrite earlier ones,
// but missing source files will not cause the removal of previously written destination files.
// The first two arguments must be paths to directories.
// The dst directory will be created, if it doesn't exist.
func StoreOrLoadFiles(srcFS, dstFS vfs.FileSystem, src, dst string, filenames ...*FilenameMapping) error {
	// verify that source path exists and is a directory
	srcFi, err := srcFS.Stat(src)
	if err != nil {
		return fmt.Errorf("error accessing source directory: %w", err)
	} else if !srcFi.IsDir() {
		return fmt.Errorf("source path '%s' is not a directory", src)
	}

	// check if destination directory exists and create it, if not
	dstFi, err := dstFS.Stat(dst)
	if err != nil {
		if vfs.IsNotExist(err) {
			if err := dstFS.MkdirAll(dst, os.ModePerm|os.ModeDir); err != nil {
				return fmt.Errorf("error creating destination directory: %w", err)
			}
		} else {
			return fmt.Errorf("error accessing destination directory: %w", err)
		}
	} else if !dstFi.IsDir() {
		return fmt.Errorf("destination path '%s' is not a directory", dst)
	}

	// copy or remove files
	writtenFiles := sets.New[string]() // keep track of already written files, because these are not removed if a later mapping specifies a non-existing source file
	for _, mapping := range filenames {
		srcPath := filepath.Join(src, mapping.Src)
		dstPath := filepath.Join(dst, mapping.Dst)

		var data []byte
		var err error
		srcIsSymlink := false
		// check if source file is a symlink
		if mapping.SymlinkPossible {
			sdata, err := srcFS.Readlink(srcPath)
			if err == nil {
				debug.Debug("Source file '%s' is a symlink", srcPath)
				srcIsSymlink = true
				data = []byte(sdata)
			}
		}

		if !srcIsSymlink {
			data, err = vfs.ReadFile(srcFS, srcPath)
		}
		if err != nil {
			if !vfs.IsNotExist(err) {
				return fmt.Errorf("error reading file '%s': %w", srcPath, err)
			}
			// remove destination file, unless we have written it before
			if writtenFiles.Has(dstPath) {
				debug.Debug("Skipping removal of destination file '%s' because it was overwritten before", dstPath)
				continue
			}
			if err := dstFS.Remove(dstPath); err != nil {
				if !vfs.IsNotExist(err) {
					return fmt.Errorf("error removing file '%s': %w", dstPath, err)
				}
				debug.Debug("Neither source '%s' nor destination '%s' file exist", srcPath, dstPath)
				continue
			}
			debug.Debug("Removed destination file '%s' because source file does not exist", dstPath)
			continue
		}

		if mapping.SymlinkPossible {
			// if the destination file could be a symlink, remove it first to ensure that we write to the actual destination path and not to a linked one
			if err := dstFS.Remove(dstPath); err != nil && !vfs.IsNotExist(err) {
				return fmt.Errorf("error removing file '%s': %w", dstPath, err)
			}
		}
		if srcIsSymlink {
			if err := dstFS.Symlink(string(data), dstPath); err != nil {
				return fmt.Errorf("error creating symlink '%s' to '%s': %w", dstPath, data, err)
			}
		} else {
			if err := vfs.WriteFile(dstFS, dstPath, data, os.ModePerm); err != nil {
				return fmt.Errorf("error writing file '%s': %w", dstPath, err)
			}
		}
		writtenFiles.Insert(dstPath)
		debug.Debug("Copied '%s' to '%s'", srcPath, dstPath)
	}

	return nil
}

// FilenameMapping is a small helper struct that represents a mapping from a source filename to a destination filename
// We use a list of these instead of a map to ensure that the order is preserved.
// If SymlinkPossible is true, the symlink handling logic is enabled.
type FilenameMapping struct {
	Src             string
	Dst             string
	SymlinkPossible bool
}

// Fnm is a helper constructor for FilenameMapping
func Fnm(src, dst string) *FilenameMapping {
	return &FilenameMapping{Src: src, Dst: dst}
}

// Fnmi is a helper constructor for FilenameMapping, that maps the source filename to itself
func Fnmi(src string) *FilenameMapping {
	return &FilenameMapping{Src: src, Dst: src}
}

func (f *FilenameMapping) WithPossibleSymlink() *FilenameMapping {
	f.SymlinkPossible = true
	return f
}
