## kw load

Load the stored configuration

### Synopsis

Load the configuration that is stored under the given id.

Simply speaking, 'kw store' stores the current kubeconfig and 'kw load' can then be used to load it again.
Storing the kubeconfig does not change it in any way, loading overwrites the current kubeconfig with the stored one.

The id must be a valid filename consisting of alphanumerical characters, starting with a letter.
If no id is given, it is defaulted to 'default'.

Note that loading kubeconfigs might not work if they have been created by plugins which have side effects when creating the kubeconfig,
as only the kubeconfig and the plugin's state are restored, but the side effects cannot be reproduced.

Loading an id that does not exist will result in an error and not change the current configuration.

```
kw load [<id>] [flags]
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

* [kw](kw.md)	 - Quickly switch between multiple Kubernetes clusters

