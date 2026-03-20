# Plugin Contract

A 'plugin' in kubeswitcher's sense is basically a binary whose behavior adheres to a specific contract, which is explained in more detail below.

The [library module](pluginlib/) offers some utility functions to make developing plugins easier. Documentation is to follow.

## Generic State

The generic kubeswitcher state stores information that is relevant independently of any used plugins.
```yaml
lastUsed: # contains information about the last used kubeswitcher command
  command: kw foo bar # the last used kubeswitcher command
  plugin: foo # the name of the plugin that handled the last command (this determines the format of the plugin-specific state)
```

## Plugin State

Each plugin may store its own state. While kubeswitcher cannot understand the semantics behind it, it will be shown on the `kw info` command and can be used by the plugin itself on repeated calls.

Structure and content depend on the plugin. A 'Gardener' plugin might for example store landscape, project, and shoot name, while a plugin that simply uses a kubeconfig from the clipboard doesn't need any state at all.

## Contract

This is the contract between the kubeswitcher tool and its plugin binaries.

The plugin binary will get the following env vars:
- `KUBESWITCHER_KUBECTL_PATH`
  - Path to the kubectl binary. Should default to `kubectl` if not set.
    - MAY be used when `kubectl` needs to be called.
    - MUST NOT be written.
- `KUBESWITCHER_KUBECONFIG_PATH`
  - Path to the kubeconfig file.
    - MAY be read.
    - MAY be written.
- `KUBESWITCHER_CURRENT_PLUGIN_NAME`
  - This is the name under which the plugin is registered in kubeswitcher.
- `KUBESWITCHER_GENERIC_STATE_PATH`
  - Path to the file containing the generic state. The file is JSON-formatted and the state has not been updated with any information from the command that is currently being executed.
    - MAY be read.
    - MUST NOT be written.
  - The file might not exist if the current command is the first kubeswitcher command to be executed. If it does not exist, the plugin state can be assumed to not exist either.
  - Mainly relevant for plugins since it contains the name of the plugin that handled the last call. If that name matches `KUBESWITCHER_CURRENT_PLUGIN_NAME`, then the plugin state comes from the same plugin and can be parsed.
- `KUBESWITCHER_PLUGIN_STATE_PATH`
  - Path to the file containing the plugin state. The file is JSON-formatted.
    - MAY be read if the value of `KUBESWITCHER_CURRENT_PLUGIN_NAME` matches `lastUsed.plugin` in the generic state.
    - MUST be modified (written or deleted) if the plugin exits with an exit code of `0` and changed the kubeconfig. MUST NOT be modified on a non-zero exit code or if the kubeconfig was not changed.
  - The plugin should only evaluate the content of this file if the current plugin name (from the `KUBESWITCHER_CURRENT_PLUGIN_NAME`) matches the last used plugin name (from `lastUsed.plugin` in the generic state).
  - The plugin must either overwrite this file with its own state after the command has been executed or delete it, in case the plugin doesn't use state.
- `KUBESWITCHER_NOTIFICATION_MESSAGE_PATH`
  - Path to a file where a 'kubeconfig has changed' notification message can be written into.
    - MUST NOT be read (doesn't exist when the plugin is called).
    - MUST be written if the plugin modified the kubeconfig. SHOULD NOT be empty.
  - If the plugin changed the kubeconfig, it must write a notification message into this file.
  - The existence of this file is used to evaluate if the kubeconfig has changed, so even if no notification message should be printed, an empty file has to be created if the kubeconfig changed.
  - Example message: `Switched kubeconfig to shoot 'my-shoot' in project 'my-project' on Gardener live landscape.`
- `KUBESWITCHER_ID_PATH`
  - Path to the file containing the id.
    - MAY be read (probably not helpful, though).
    - MUST be written if the plugin modified the kubeconfig. SHOULD NOT be empty.
  - The contents of this file can be used to display the current kubeconfig target as part of a shell prompt.
    - Should be as short as possible while still containing precise information about the currently targeted cluster.
    - Best practice: Prefix it with the value from `KUBESWITCHER_CURRENT_PLUGIN_NAME`. Also, don't add line breaks.
    - Example: If the Gardener plugin was used, the id could look like this: `garden:<landscape>/<project>/<shoot>`.
  - This is also used by the command which shows the history.
- `KUBESWITCHER_INTERNAL_CALL_PATH`
  - Path to the file for command referrence.
    - MUST NOT be read (doesn't exist when the plugin is called).
    - MAY be written.
  - If the plugin wraps another kubeswitcher command, that one cannot be called directly (calling the binary would mess with state handling). Instead, the calling command can write the subcommand and all arguments to this file path.
    - Example: If your command wants to call `kw exec foo --bar`, with `exec` being the subcommand, it should write `exec foo --bar` into the file.
    - The file must contain only a single command.
  - If the plugin command needs to be called again after the internal call has been executed (to react on it or modify the state), it has to write a callback file, see below.
  - Internal calls can cause further internal calls.
  - Note that, because plugin subcommand names are configurable, in order for plugin A to call plugin B, the subcommand name for B has to be part of A's configuration.
    - The built-in subcommands are named `<internal_X>`, where `X` is the subcommand.
- `KUBESWITCHER_INTERNAL_CALLBACK_REQUEST_PATH`
  - Path to the file to request a callback after the internal call.
    - MAY be written.
    - MUST NOT be read (doesn't exist when the plugin is called).
  - Some commands that wrap other commands might be fine with just calling that other command, while others want to react on the command's result or modify the kubeswitcher state afterwards. For the latter case, a callback can be requested.
    - If a command requests an internal call and also creates a file with arbitrary content at this path, it will be called again after the internal call has been resolved.
      - During the callback, the information that was written to this path can be retrieved from `KUBESWITCHER_INTERNAL_CALLBACK_STATE_PATH`.
    - The command can use this file to store arbitrary information that it might need to resolve the later callback.
- `KUBESWITCHER_INTERNAL_CALLBACK_STATE_PATH`
  - Path to the file for a callback after the internal call.
    - MUST NOT be written.
    - MUST be read, but only if the command makes use of internal calls with callbacks. Can be ignored otherwise.
  - If a command requested an internal call and also created a file with arbitrary content at `KUBESWITCHER_INTERNAL_CALLBACK_REQUEST_PATH`, it will be called again after the internal call has been resolved, with the same content being available at this path.
    - This means that these commands have to differentiate between two cases:
      - Called and no file exists at this path: Regular command call.
      - Called and the file exists: Callback after an internal call.
- `KUBESWITCHER_PLUGIN_CONFIG`
  - Contains the static plugin config as JSON, if specified. Is unset otherwise.
- `KUBESWITCHER_SESSION_ID`
  - Contains the current session id, if the plugin needs this for whatever reason.
  - The session id is the same for command calls issued from the same terminal.
- `KUBESWITCHER_SESSION_CONFIG_DIR`
  - Contains the path to the session config directory.
  - This is usually the directory in which the kubeconfig and state files are stored.
  - It's usually a temporary directory and a different one exists for each terminal from which kubeswitcher is called. If you want to store permanent and global information, use the directory from `KUBESWITCHER_CONFIG_DIR` instead.
  - Can be used if the plugin wants to store information in addition to its state.
  - Best practise: don't write directly into this directory, but create a `plugin_<plugin-name>` subfolder to avoid filename conflicts with kubeswitcher or other plugins
- `KUBESWITCHER_CONFIG_DIR`
  - Contains the path to the global kubeswitcher config directory.
  - Opposed to the session directory, the files from this directory are shared between all terminal processes and not lost when the terminal is closed or the system is shut down.
  - Best practise: don't write directly into this directory, but create a `plugin_<plugin-name>` subfolder to avoid filename conflicts with kubeswitcher or other plugins
- `KUBESWITCHER_FLAG_DEBUG`
  - If set and has value `true`, the `--debug` flag has been set for kubeswitcher and the plugin should print debug statements to stderr, if possible.


## Plugin Configuration

Plugins have to be configured in the config file for kubeswitcher to recognize them.
```yaml
plugins:
- name: mcp # name of the subcommand that will call the plugin
  aliases: # list of aliases for the subcommand
  - m
  binary: /usr/local/bin/kw_mcp # path to the binary
  short: Switches between clusters of an MCP landscape. # short description for showing in the help
  config: <...> # will be contained in KUBESWITCHER_PLUGIN_CONFIG as json when the plugin is called
  env: # additional environment variables that should be set when calling the binary
    MY_CUSTOM_ENV: my-custom-value
```

The relevant fields of a plugin configuration are:
- `name` _required_
  - This is the name of the plugin and also the name of the subcommand under which the plugin will be registered.
    - In the above example, `kw mcp` would result in a call of the registered plugin.
  - Must be unique.
- `aliases` _optional_
  - A list of aliases for the subcommand.
  - Take care to not clash with the names or aliases of other plugins or built-in commands.
- `binary` _required_
  - Path to the binary for the command.
  - Must either be an absolute path or a single filename (which is then resolved via the `PATH` env var).
- `short` _required_
  - This should be a very short description of the plugin, which is shown next to the plugin's name when `kw --help` is used.
- `config` _optional_
  - Arbitrary data that will be parsed into json and stored in the `KUBESWITCHER_PLUGIN_CONFIG` env var when the plugin is called.
- `env` _optional_
  - Statically defined environment variables that will be set when the plugin is called.
