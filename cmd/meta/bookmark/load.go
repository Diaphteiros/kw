package bookmark

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/storage"
	"github.com/Diaphteiros/kw/pluginlib/pkg/selector"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

const IdRegexString = `^[a-zA-Z]+(?:-*[0-9a-zA-Z]+)*$`

var IdRegex = regexp.MustCompile(IdRegexString)

var LoadCmd = &cobra.Command{
	Use:     "load [<key>]",
	Aliases: []string{"l"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Load the bookmarked configuration",
	Long: `Load the configuration that is bookmarked under the given key.

Simply speaking, the 'save' subcommand stores the current kubeconfig and this one can then be used to load it again.
Bookmarking the kubeconfig does not change it in any way, loading overwrites the current kubeconfig with the stored one.

The key must consist of alphanumerical characters and dashes only, and it must neither begin nor end with a dash.
If no key is given, you will be prompted for one.

Note that loading kubeconfigs might not work if they have been created by plugins which have side effects when creating the kubeconfig,
as only the kubeconfig and the plugin's state are restored, but the side effects cannot be reproduced.

Loading a key that does not exist will result in an error and not change the current configuration.

Note that the bookmarks are shared between all terminal sessions.`,
	Run: func(cmd *cobra.Command, args []string) {
		var arg string
		if len(args) > 0 {
			arg = args[0]
		}
		key := bookmarkIndexFromArgumentOrPrompt(arg, nil)
		validateStoreLoadKey(key)
		if err := storage.ManualLoad(key); err != nil {
			libutils.Fatal(1, "error loading bookmark: %w\n", err)
		}
	},
}

func validateStoreLoadKey(key string) {
	if !IdRegex.MatchString(key) {
		libutils.Fatal(1, "key '%s' is not a valid id (regex is '%s')\n", key, IdRegexString)
	}
}

// bookmarkIndexFromArgumentOrPrompt takes an optional argument and returns a bookmark key.
// If the argument matches a bookmark key, it is returned as is.
// Otherwise, the user is prompted with a fuzzy finder to select an existing bookmark key.
// This function will fatal if no bookmark key can be selected.
// If the given bookmarkInventory is nil, it will be loaded internally.
func bookmarkIndexFromArgumentOrPrompt(key string, si storage.StorageInventory) string {
	if si == nil {
		var err error
		si, err = storage.RenderStorageInventory()
		if err != nil {
			libutils.Fatal(1, "error creating internal bookmark representation: %w\n", err)
		}
	}
	idx := -1
	if key != "" {
		// if the user entered an argument, check if it matches an existing key
		for i, se := range si {
			if se.Key == key {
				idx = i
				break
			}
		}
	}
	if idx < 0 {
		// the user either did not enter a key or it did not match any existing key
		// use the fuzzy finder to select one
		_, elem, err := selector.New[*storage.StorageEntry]().
			WithFatalOnAbort("no bookmark key selected\n").
			WithFatalOnError("error while selecting bookmark key: %w\n").
			From(si, func(se *storage.StorageEntry) string {
				return se.Key
			}).
			WithPreview(func(se *storage.StorageEntry, _, _ int) string {
				return fmt.Sprintf("Bookmark Key: %s\nIdentity: %s", se.Key, se.Id)
			}).
			WithPrompt("Choose a bookmark key: ").
			WithQuery(key).
			Select()
		if err != nil {
			libutils.Fatal(1, "error selecting bookmark key: %w\n", err)
		}
		key = elem.Key
	}
	return key
}
