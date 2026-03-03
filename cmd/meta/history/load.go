package history

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pkg/storage"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	"github.com/Diaphteiros/kw/pluginlib/pkg/selector"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var HistoryLoadCmd = &cobra.Command{
	Use:     "load [<index>]",
	Aliases: []string{"l"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Load a configuration from history",
	Long: `Load a configuration from history.
	
This command requires a history depth of at least 1.
The index must be an integer between 0 (inclusive) and the history depth (exclusive).
Basically, 1 refers to the last configuration before the current one, 2 to the one before that, and so on.
This means that (<history depth> - 1) refers to the oldest configuration in the history.

Use the 'history view' command to see the available indices.

'history load 0' is a no-op, since it refers to the current configuration.

If no index is specified, the history is shown and the user is prompted for an index to load.`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Runtime.Config().Kubeswitcher.HistoryDepth == 0 {
			libutils.Fatal(1, "history is disabled\n")
		}
		var idx int
		var err error
		if len(args) > 0 {
			idx, err = strconv.Atoi(args[0])
			if err != nil {
				libutils.Fatal(1, "index argument must be an integer, got '%s': %w\n", args[0], err)
			}
		} else {
			// load history
			hist, err := storage.RenderHistory(Global)
			if err != nil {
				libutils.Fatal(1, "error creating internal history representation: %w\n", err)
			}
			idx = promptForHistoryIndex(hist)
		}
		if idx < 0 || idx >= config.Runtime.Config().Kubeswitcher.HistoryDepth {
			libutils.Fatal(1, "index must be between 0 (inclusive) and %d (exclusive), got %d\n", config.Runtime.Config().Kubeswitcher.HistoryDepth, idx)
		}
		currentHistoryIndex, err := storage.GetCurrentHistoryIndex(Global)
		if err != nil {
			libutils.Fatal(1, "error getting current history index: %w\n", err)
		}
		if currentHistoryIndex == -1 {
			libutils.Fatal(1, "history is empty or current history index was lost\n")
		}
		convertedIdx := storage.ConvertHistoryIndexAndDirName(currentHistoryIndex, idx)
		if convertedIdx < 0 || convertedIdx >= config.Runtime.Config().Kubeswitcher.HistoryDepth {
			libutils.Fatal(1, "index %d does not seem to match a valid history entry\n", idx)
		}
		if err := storage.HistoryLoad(fs.FS, fs.FS, strconv.Itoa(convertedIdx), Global); err != nil {
			libutils.Fatal(1, "error loading history entry: %w\n", err)
		}
	},
}

func promptForHistoryIndex(hist storage.HistoryInventory) int {
	_, elem, err := selector.New[*storage.HistoryEntry]().
		WithPrompt("Choose a history index: ").
		WithFatalOnAbort("no history index selected").
		WithFatalOnError("error while selecting history index: %w\n").
		From(hist, func(he *storage.HistoryEntry) string {
			return he.String()
		}).
		Select()
	if err != nil {
		libutils.Fatal(1, "error while selecting history index: %w\n", err)
	}
	return elem.Key
}
