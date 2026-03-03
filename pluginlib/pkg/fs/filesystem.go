package fs

import (
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
)

// FS is the default filesystem for all file operations.
// Usually, this defaults to the OS filesystem.
// For tests, this can be replaced with a virtual filesystem.
var FS vfs.FileSystem

// MFS is a virtual, in-memory filesystem that can be used for file operations that don't need to be persisted.
var MFS vfs.FileSystem

func init() {
	FS = osfs.New()
	MFS = memoryfs.New()
}
