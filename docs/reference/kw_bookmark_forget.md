## kw bookmark forget

Delete a bookmark entry

### Synopsis

Delete a bookmark entry.
	
Either one or more bookmark keys or the '--all' flag must be specified.
If bookmark keys are specified, the corresponding bookmark entries will be deleted.
Missing entries are ignored.

If the '--all' flag is set, all bookmark entries will be deleted.

Note that the bookmarks are shared between all terminal sessions.

```
kw bookmark forget [<name> ...] [flags]
```

### Options

```
      --all    Delete all bookmark entries
  -h, --help   help for forget
```

### Options inherited from parent commands

```
      --debug   Print debug information to stderr.
```

### SEE ALSO

* [kw bookmark](kw_bookmark.md)	 - Interact with the kubeconfig bookmarks

