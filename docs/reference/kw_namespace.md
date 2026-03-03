## kw namespace

Change the default namespace in the current context of the kubeconfig

### Synopsis

Change the default namespace in the current context of the kubeconfig.

This is basically the same thing as running 'kubectl config set-context --current --namespace=<namespace>'.
If called without any argument, the command fetches the namespaces from the currently selected cluster and prompts for a selection.

Note that this command does change the kubeconfig file, but doesn't create a new kubeswitcher history entry.

```
kw namespace [<namespace>] [flags]
```

### Options

```
  -h, --help   help for namespace
```

### Options inherited from parent commands

```
      --debug   Print debug information to stderr.
```

### SEE ALSO

* [kw](kw.md)	 - Quickly switch between multiple Kubernetes clusters

