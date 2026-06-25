## kw config view

View the kubeswitcher configuration file

### Synopsis

View the kubeswitcher configuration file.

Use the '--core' flag to exclude plugin configuration from the output.
To only view a specific plugin's configuration, use the '--plugin <plugin_name>' flag.

Usually, the configuration file is parsed and before being printed. This means that comments don't appear in the output and some fields might have changed due to defaulting.
To view the raw configuration file, use the '--raw' flag.

```
kw config view [flags]
```

### Options

```
  -c, --core            Print only the core configuration, excluding plugin configurations.
  -h, --help            help for view
  -o, --output string   Output format. Valid formats are [json, yaml]. (default "yaml")
  -p, --plugin string   Print only the configuration for the specified plugin.
  -r, --raw             Print the raw configuration file without parsing it. Cannot be combined with any arguments that filter the output or change its format.
```

### Options inherited from parent commands

```
      --debug   Print debug information to stderr.
```

### SEE ALSO

* [kw config](kw_config.md)	 - Interact with the kubeswitcher configuration

