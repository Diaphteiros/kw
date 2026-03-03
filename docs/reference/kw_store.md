## kw store

Store the current configuration

### Synopsis

Store the current configuration under the given id.

Simply speaking, 'kw store' stores the current kubeconfig and 'kw load' can then be used to load it again.
Storing the kubeconfig does not change it in any way, loading overwrites the current kubeconfig with the stored one.
Subsequent calls with the same id will overwrite the previously stored configuration with the current one.

The id must be a valid filename consisting of alphanumerical characters, starting with a letter.
If no id is given, the kubeconfig is stored under the default id 'default'.

Note that loading kubeconfigs might not work if they have been created by plugins which have side effects when creating the kubeconfig,
as only the kubeconfig and the plugin's state are restored, but the side effects cannot be reproduced.

Stored configurations are never deleted, so if your session directory is not a temporary one, you might want to clean it up from time to time.

```
kw store [<id>] [flags]
```

### Options

```
  -h, --help   help for store
```

### Options inherited from parent commands

```
      --debug   Print debug information to stderr.
```

### SEE ALSO

* [kw](kw.md)	 - Quickly switch between multiple Kubernetes clusters

