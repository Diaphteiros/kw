## kw history load

Load a configuration from history

### Synopsis

Load a configuration from history.
	
This command requires a history depth of at least 1.
The index must be an integer between 0 (inclusive) and the history depth (exclusive).
Basically, 1 refers to the last configuration before the current one, 2 to the one before that, and so on.
This means that (<history depth> - 1) refers to the oldest configuration in the history.

Use the 'history view' command to see the available indices.

'history load 0' is a no-op, since it refers to the current configuration.

If no index is specified, the history is shown and the user is prompted for an index to load.

```
kw history load [<index>] [flags]
```

### Options

```
  -h, --help   help for load
```

### Options inherited from parent commands

```
      --debug    Print debug information to stderr.
  -g, --global   Use the global history instead of the session-specific one
```

### SEE ALSO

* [kw history](kw_history.md)	 - Interact with the history

