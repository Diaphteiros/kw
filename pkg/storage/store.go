package storage

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
)

// ManualStore handles a store operation that was caused by a manual call to the store command
func ManualStore(id string) error {
	debug.Debug("Storing to storage with id '%s'", id)
	if len(id) == 0 {
		return fmt.Errorf("store id must not be empty")
	}
	spath := GetStoragePath(id)
	if err := StoreOrLoadFiles(fs.FS, fs.FS, config.Runtime.SessionDir(), spath, allFileMappings...); err != nil {
		return fmt.Errorf("error storing configuration: %w", err)
	}
	return nil
}

// HistoryStore stores the current configuration in the history.
// This is only triggered by internal calls, not by manual calls to the store command.
func HistoryStore(srcFS, dstFS vfs.FileSystem, id string) error {
	debug.Debug("Storing to history with id '%s'", id)
	if config.Runtime.Config().Kubeswitcher.HistoryDepth <= 0 {
		debug.Debug("Skipping store because history is disabled")
		return nil
	}
	if len(id) == 0 {
		return fmt.Errorf("history store id must not be empty")
	}
	spath := GetHistoryPath(id, false)
	if err := StoreOrLoadFiles(srcFS, dstFS, config.Runtime.SessionDir(), spath, allFileMappings...); err != nil {
		return fmt.Errorf("error storing history: %w", err)
	}
	return nil
}

// StoreFromLocalToGlobalHistory stores current local history entry as newest global history entry.
// This is done by creating a symlink in the global history directory that points to the directory in the local history.
func StoreFromLocalToGlobalHistory() error {
	debug.Debug("Storing local history entry to global history")
	if config.Runtime.Config().Kubeswitcher.HistoryDepth <= 0 {
		debug.Debug("Skipping global store because history is disabled")
		return nil
	}
	globalId, err := GetCurrentHistoryIndex(true)
	if err != nil {
		return err
	}
	globalId = (globalId + 1) % config.Runtime.Config().Kubeswitcher.HistoryDepth
	gPath := GetHistoryPath(strconv.Itoa(globalId), true)
	if err := fs.FS.MkdirAll(GetHistoryRootPath(true), os.ModePerm|os.ModeDir); err != nil {
		return fmt.Errorf("error creating global history directory '%s': %w", gPath, err)
	}
	debug.Debug("Renaming last global history entry to '%s'", gPath)
	old, err := fs.FS.Lstat(gPath)
	if err == nil {
		if old.IsDir() {
			if err := fs.FS.RemoveAll(gPath); err != nil {
				return fmt.Errorf("error removing existing directory at '%s': %w", gPath, err)
			}
		} else {
			if err := fs.FS.Remove(gPath); err != nil {
				return fmt.Errorf("error removing existing file at '%s': %w", gPath, err)
			}
		}
	} else if !vfs.IsNotExist(err) {
		return fmt.Errorf("error accessing existing global history entry '%s': %w", gPath, err)
	}
	// old file is removed, now create symlink
	localId, err := GetCurrentHistoryIndex(false)
	if err != nil {
		return err
	}
	lPath := GetHistoryPath(strconv.Itoa(localId), false)
	debug.Debug("Creating symlink from '%s' to '%s'", gPath, lPath)
	if err := fs.FS.Symlink(lPath, gPath); err != nil {
		return fmt.Errorf("error creating symlink from '%s' to '%s': %w", gPath, lPath, err)
	}
	if err := vfs.WriteFile(fs.FS, GetCurrentHistoryIndexFilePath(true), []byte(strconv.Itoa(globalId)), os.ModePerm); err != nil {
		return fmt.Errorf("error writing global history index: %w", err)
	}
	return nil
}

// StoreToTmpHistory stores the current configuration in the temporary history.
// This is meant to be called before a command that might change the configuration is executed.
func StoreToTmpHistory() error {
	debug.Debug("Storing to temporary history (from regular filesystem to memory filesystem)")
	if config.Runtime.Config().Kubeswitcher.HistoryDepth <= 0 {
		debug.Debug("Skipping store because history is disabled")
		return nil
	}
	return HistoryStore(fs.FS, fs.MFS, HistoryTmpKey)
}

// StoreFromCurrentToHistory copies the current state to the history.
// This is meant to be called after a command that has changed the configuration has been executed.
// This updates the current history index.
func StoreFromCurrentToHistory() error {
	debug.Debug("Storing from current state to history")
	if config.Runtime.Config().Kubeswitcher.HistoryDepth <= 0 {
		debug.Debug("Skipping store because history is disabled")
		return nil
	}
	curPath := config.Runtime.SessionDir()
	currentHistoryIndex, err := GetCurrentHistoryIndex(false)
	if err != nil {
		return err
	}
	currentHistoryIndex = (currentHistoryIndex + 1) % config.Runtime.Config().Kubeswitcher.HistoryDepth
	historyPath := GetHistoryPath(strconv.Itoa(currentHistoryIndex), false)
	if err := StoreOrLoadFiles(fs.FS, fs.FS, curPath, historyPath, allFileMappings...); err != nil {
		return fmt.Errorf("error storing history: %w", err)
	}
	if err := vfs.WriteFile(fs.FS, GetCurrentHistoryIndexFilePath(false), []byte(strconv.Itoa(currentHistoryIndex)), os.ModePerm); err != nil {
		return fmt.Errorf("error writing current history index: %w", err)
	}
	return nil
}
