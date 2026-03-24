package bookmark

import (
	"bufio"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/pkg/storage"

	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var SaveCmd = &cobra.Command{
	Use:     "save [<key>]",
	Aliases: []string{"store", "s"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Bookmark the current configuration",
	Long: `Bookmark the current configuration under the given key.

Simply speaking, this command bookmarks the current kubeconfig and the 'load' subcommand can then be used to load it again.
Bookmarking the kubeconfig does not change it in any way, loading overwrites the current kubeconfig with the stored one.
Subsequent calls with the same key will overwrite the previously stored configuration with the current one.

The key must consist of alphanumerical characters and dashes only, and it must neither begin nor end with a dash.
If no key is given, you will be prompted for one.

Note that loading kubeconfigs might not work if they have been created by plugins which have side effects when creating the kubeconfig,
as only the kubeconfig and the plugin's state are restored, but the side effects cannot be reproduced.

Note that the bookmarks are shared between all terminal sessions.`,
	Run: func(cmd *cobra.Command, args []string) {
		var key string
		if len(args) > 0 {
			key = args[0]
		} else {
			key = promptForBookmarkKey(cmd)
		}
		validateStoreLoadKey(key)
		if err := storage.ManualStore(key); err != nil {
			libutils.Fatal(1, "error storing configuration: %w\n", err)
		}
	},
}

func promptForBookmarkKey(cmd *cobra.Command) string {
	prompt := "Choose a bookmark key: "
	cmd.Print(prompt)
	r := bufio.NewReader(cmd.InOrStdin())
	input, err := r.ReadString('\n')
	if err != nil {
		libutils.Fatal(1, "error while reading user input: %w\n", err)
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}
