## kw

Quickly switch between multiple Kubernetes clusters

### Synopsis

Quickly switch between multiple Kubernetes clusters.

This tool allows to switch between different kubeconfig files efficiently.

There are a few built-in commands - e.g. for switching to a specific kubeconfig file or to the kubeconfig currently contained in the clipboard - and the possibility to add custom commands via plugins.
The subcommands in the 'meta' group allow switch back to recently used kubeconfigs (via the history) or to explicitly stored configurations.

A configuration file is required to use this tool and to register plugins. It is created automatically if it doesn't exist. See the 'config' subcommand for more information.

If 'TERM_SESSION_ID' is not set, a session id must be provided by setting 'KW_SESSION_ID' to a UUID-like string. This is used to create a session-specific temporary directory for storing the kubeconfig and the tool's state.
This means that each terminal session has its own state and kubeconfig, unless it shares the session id with another session.

Note that this tool expects to 'own' the KUBECONFIG environment variable and will print a warning if it doesn't point to the expected path.
The 'kubectl-env' subcommand can help with setting the KUBECONFIG environment variable to the correct path.

Using the tool with a different KUBECONFIG path than the one it manages is not recommended and might lead to unexpected behavior.
It is strongly discouraged to modify the kubeconfig that is managed by this tool by any other means than this tool itself.


### Options

```
      --debug   Print debug information to stderr.
  -h, --help    help for kw
```

### SEE ALSO

* [kw bookmark](kw_bookmark.md)	 - Interact with the kubeconfig bookmarks
* [kw config](kw_config.md)	 - Interact with the kubeswitcher configuration
* [kw custom](kw_custom.md)	 - Switch to the kubeconfig at the specified path
* [kw flip](kw_flip.md)	 - Flip the current configuration with the previously used one
* [kw history](kw_history.md)	 - Interact with the history
* [kw info](kw_info.md)	 - Shows information about the current configuration
* [kw kubectl-env](kw_kubectl-env.md)	 - Generate a script that points the KUBECONFIG env var to the kubeconfig for the current kw session
* [kw namespace](kw_namespace.md)	 - Change the default namespace in the current context of the kubeconfig
* [kw prompt](kw_prompt.md)	 - Generate a script that generates a prompt to display in the shell
* [kw repeat](kw_repeat.md)	 - Switch to the last used configuration (cross-session)
* [kw temporary](kw_temporary.md)	 - Switch to the kubeconfig currently contained in the clipboard
* [kw version](kw_version.md)	 - Print the version

