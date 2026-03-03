## kw bookmark load

Load the bookmarked configuration

### Synopsis

Load the configuration that is bookmarked under the given key.

Simply speaking, the 'save' subcommand stores the current kubeconfig and this one can then be used to load it again.
Bookmarking the kubeconfig does not change it in any way, loading overwrites the current kubeconfig with the stored one.

The key must consist of alphanumerical characters and dashes only, and it must neither begin nor end with a dash.
If no key is given, you will be prompted for one.

Note that loading kubeconfigs might not work if they have been created by plugins which have side effects when creating the kubeconfig,
as only the kubeconfig and the plugin's state are restored, but the side effects cannot be reproduced.

Loading a key that does not exist will result in an error and not change the current configuration.

Note that the bookmarks are shared between all terminal sessions.

```
kw bookmark load [<key>] [flags]
```

### Options

```
  -h, --help   help for load
```

### Options inherited from parent commands

```
      --debug   Print debug information to stderr.
```

### SEE ALSO

* [kw bookmark](kw_bookmark.md)	 - Interact with the kubeconfig bookmarks

