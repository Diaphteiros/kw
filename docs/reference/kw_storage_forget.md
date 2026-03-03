## kw storage forget

Delete a storage entry

### Synopsis

Delete a storage entry.
	
Either one or more storage keys or the '--all' flag must be specified.
If storage keys are specified, the corresponding storage entries will be deleted.
Missing entries are ignored.

If the '--all' flag is set, all storage entries will be deleted.

Note that the storage is shared between all terminal sessions.

```
kw storage forget [<name> ...] [flags]
```

### Options

```
      --all    Delete all storage entries
  -h, --help   help for forget
```

### Options inherited from parent commands

```
      --debug   Print debug information to stderr.
```

### SEE ALSO

* [kw storage](kw_storage.md)	 - Interact with the kubeconfig storage

