## kw config

Interact with the kubeswitcher configuration

### Synopsis

Interact with the kubeswitcher configuration.

The default configuration directory is $HOME/.kubeswitcher_config/.
This can be changed by setting the environment variable KW_CONFIG_REPO.

Within the configuration directory, a 'config.yaml' file is expected to contain the kubeswitcher configuration.
Directory and file (with default values) will be created when any kubeswitcher command is run, if they do not exist.

### Options

```
  -h, --help   help for config
```

### Options inherited from parent commands

```
      --debug   Print debug information to stderr.
```

### SEE ALSO

* [kw](kw.md)	 - Quickly switch between multiple Kubernetes clusters
* [kw config path](kw_config_path.md)	 - View the path of the kubeswitcher configuration file

