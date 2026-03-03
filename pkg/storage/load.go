package storage

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pkg/utils"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

// ManualLoad handles a load operation that was caused by a manual call to the load command.
func ManualLoad(id string) error {
	debug.Debug("Loading from storage with id '%s'", id)
	if len(id) == 0 {
		return fmt.Errorf("load id must not be empty")
	}
	if _, err := fs.FS.Stat(GetStoragePath(id)); err != nil && vfs.IsNotExist(err) {
		return fmt.Errorf("id '%s' does not exist", id)
	}
	if err := StoreOrLoadFiles(fs.FS, fs.FS, GetStoragePath(id), config.Runtime.SessionDir(), allFileMappings...); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}
	return nil
}

// HistoryLoad loads the configuration from the history.
func HistoryLoad(srcFS, dstFS vfs.FileSystem, id string, global bool) error {
	debug.Debug("Loading from history (global: %t) with id '%s'", global, id)
	if len(id) == 0 {
		return fmt.Errorf("history load id must not be empty")
	}
	if err := StoreOrLoadFiles(srcFS, dstFS, GetHistoryPath(id, global), config.Runtime.SessionDir(), allFileMappings...); err != nil {
		return fmt.Errorf("error loading history: %w", err)
	}
	return nil
}

// LoadFromTmpHistory loads the configuration from the temporary history.
// This can be used to recover the original configuration after a command has failed.
func LoadFromTmpHistory() error {
	debug.Debug("Loading from temporary history (from memory filesystem to regular filesystem)")
	return HistoryLoad(fs.MFS, fs.FS, HistoryTmpKey, false)
}

// HistoryInventory is an internal representation of the history.
// The file-system history is converted into this for better internal handling and display.
type HistoryInventory []*HistoryEntry
type HistoryEntry abstractInventoryEntry[int]

var _ json.Marshaler = HistoryInventory{}

// String returns a human-readable representation of the history inventory.
// Reverses the order to print the most recent entries last.
func (hi HistoryInventory) String() string {
	if len(hi) == 0 {
		return "History is empty."
	}
	tmp := hi.ToStrings()
	slices.Reverse(tmp)
	return strings.Join(tmp, "\n")
}

func (he *HistoryEntry) String() string {
	return ((*abstractInventoryEntry[int])(he)).String()
}

// ToStrings converts the history inventory into a slice of strings for easier display.
func (hi HistoryInventory) ToStrings() []string {
	return libutils.Project(hi, func(he *HistoryEntry) string {
		return he.String()
	})
}

// MarshalJSON implements json.Marshaler.
func (h HistoryInventory) MarshalJSON() ([]byte, error) {
	// ensure that index matches the actual position in the slice
	maxIdx := -1
	for _, he := range h {
		if he.Key > maxIdx {
			maxIdx = he.Key
		}
	}
	aligned := make([]*HistoryEntry, maxIdx+1)
	for _, he := range h {
		aligned[he.Key] = he
	}
	return json.Marshal(aligned)
}

type StorageInventory []*StorageEntry
type StorageEntry abstractInventoryEntry[string]

var _ json.Marshaler = StorageInventory{}

func (si StorageInventory) String() string {
	if len(si) == 0 {
		return "Storage is empty."
	}
	return strings.Join(libutils.Project(si, func(se *StorageEntry) string {
		return se.String()
	}), "\n")
}

func (se *StorageEntry) String() string {
	return ((*abstractInventoryEntry[string])(se)).String()
}

// MarshalJSON implements json.Marshaler.
func (si StorageInventory) MarshalJSON() ([]byte, error) {
	// convert to a map with index as key
	conv := map[string]*StorageEntry{}
	for _, se := range si {
		conv[se.Key] = se
	}
	return json.Marshal(conv)
}

type abstractInventoryEntry[T any] struct {
	// Key is the identifying property of the entry.
	Key T `json:"key"`
	// Id is the id (in the kubeswitcher sense of 'id') of the entry.
	Id string `json:"id"`
	// Path is the path on the filesystem where the entry is stored.
	Path string `json:"path"`
}

func (e *abstractInventoryEntry[T]) String() string {
	return fmt.Sprintf("%v: %s", e.Key, e.Id)
}

// RenderHistory converts the file-based history into an internal representation.
func RenderHistory(global bool) (HistoryInventory, error) {
	hist := make(HistoryInventory, 0, config.Runtime.Config().Kubeswitcher.HistoryDepth+1)
	hRoot := GetHistoryRootPath(global)
	hDirs, err := vfs.ReadDir(fs.FS, hRoot)
	if err != nil {
		if vfs.IsNotExist(err) {
			return nil, nil
		}
		libutils.Fatal(1, "error reading history directory: %w\n", err)
	}
	// get current history index
	currentHistoryIndex, err := GetCurrentHistoryIndex(global)
	if err != nil {
		libutils.Fatal(1, "error getting current history index: %w\n", err)
	}

	// read history entries
	for _, hDir := range hDirs {
		hDirPath := filepath.Join(hRoot, hDir.Name())
		if global {
			realDir, err := fs.FS.Readlink(filepath.Join(hRoot, hDir.Name()))
			if err == nil {
				hDirPath = realDir
			} else if !hDir.IsDir() {
				debug.Debug("Ignoring history entry '%s' because it is neither a directory nor a symlink", hDirPath)
				continue
			}
		} else {
			if !hDir.IsDir() {
				debug.Debug("Ignoring history entry '%s' because it is not a directory", hDirPath)
				continue
			}
		}
		dirIdx, err := strconv.Atoi(hDir.Name())
		if err != nil {
			continue
		}
		idx := ConvertHistoryIndexAndDirName(currentHistoryIndex, dirIdx)
		if idx < 0 || idx >= config.Runtime.Config().Kubeswitcher.HistoryDepth {
			continue
		}
		id, err := utils.GetId(hDirPath)
		if err != nil {
			libutils.Fatal(1, "error getting id for history entry: %w\n", err)
		}
		hist = append(hist, &HistoryEntry{Key: idx, Id: id, Path: hDirPath})
	}

	// sort history by index
	slices.SortFunc(hist, func(a, b *HistoryEntry) int {
		return a.Key - b.Key
	})

	return hist, nil
}

// RenderStorageInventory converts the file-based storage into an internal representation.
func RenderStorageInventory() (StorageInventory, error) {
	si := StorageInventory{}
	sRoot := GetStorageRootPath()
	sDirs, err := vfs.ReadDir(fs.FS, sRoot)
	if err != nil {
		if vfs.IsNotExist(err) {
			return nil, nil
		}
		libutils.Fatal(1, "error reading storage directory: %w\n", err)
	}
	for _, sDir := range sDirs {
		if !sDir.IsDir() {
			continue
		}
		sPath := filepath.Join(sRoot, sDir.Name())
		id, err := utils.GetId(sPath)
		if err != nil {
			libutils.Fatal(1, "error getting id for storage entry: %w\n", err)
		}
		si = append(si, &StorageEntry{Key: sDir.Name(), Id: id, Path: sPath})
	}
	return si, nil
}
