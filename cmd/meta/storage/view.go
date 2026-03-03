package storage

import (
	"github.com/spf13/cobra"
	"github.com/Diaphteiros/kw/cmd/meta/bookmark"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var output libutils.OutputFormat

var ViewCmd = &cobra.Command{
	Use:        "view",
	Aliases:    []string{"list", "ls", "v"},
	Args:       cobra.NoArgs,
	Short:      "[DEPRECATED] View the storage",
	Long:       `View the storage.`,
	Deprecated: "use 'bookmark view' instead.",
	Run:        bookmark.ViewCmd.Run,
}

func init() {
	libutils.AddOutputFlag(ViewCmd.Flags(), &output, libutils.OUTPUT_TEXT)
}
