package storage

import (
	"regexp"

	"github.com/spf13/cobra"

	"github.com/Diaphteiros/kw/cmd/meta/bookmark"
)

const IdRegexString = `^[a-zA-Z]+(?:-*[0-9a-zA-Z]+)*$`

var IdRegex = regexp.MustCompile(IdRegexString)

var LoadCmd = &cobra.Command{
	Use:     "load [<key>]",
	Aliases: []string{"l"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "[DEPRECATED] Load the stored configuration",
	Long: `Load the configuration that is stored under the given key.

Simply speaking, the 'store' subcommand stores the current kubeconfig and this one can then be used to load it again.
Storing the kubeconfig does not change it in any way, loading overwrites the current kubeconfig with the stored one.

The key must consist of alphanumerical characters and dashes only, and it must neither begin nor end with a dash.
If no key is given, you will be prompted for one.

Note that loading kubeconfigs might not work if they have been created by plugins which have side effects when creating the kubeconfig,
as only the kubeconfig and the plugin's state are restored, but the side effects cannot be reproduced.

Loading an key that does not exist will result in an error and not change the current configuration.

Note that the storage is shared between all terminal sessions.`,
	Deprecated: "use 'bookmark load' instead.",
	Run:        bookmark.LoadCmd.Run,
}
