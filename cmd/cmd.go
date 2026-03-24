package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"

	basiccmd "github.com/Diaphteiros/kw/cmd/basic"
	configcmd "github.com/Diaphteiros/kw/cmd/config"
	metacmd "github.com/Diaphteiros/kw/cmd/meta"
	misccmd "github.com/Diaphteiros/kw/cmd/misc"
	versioncmd "github.com/Diaphteiros/kw/cmd/version"
	"github.com/Diaphteiros/kw/pkg/cmdgroups"
	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pkg/storage"

	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	"github.com/Diaphteiros/kw/pluginlib/pkg/state"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

var (
	pluginStateLastModified time.Time
	builtinSubcommands      = []*cobra.Command{
		basiccmd.TemporaryCmd,
		basiccmd.CustomCmd,
		metacmd.StorageCmd,
		metacmd.BookmarkCmd,
		metacmd.FlipCmd,
		metacmd.RepeatCmd,
		metacmd.InfoCmd,
		metacmd.HistoryCmd,
		metacmd.NamespaceCmd,
		configcmd.ConfigCmd,
		misccmd.KubectlEnvCmd,
		misccmd.PromptCmd,
		versioncmd.VersionCmd,
	}
	skipKubeconfigWarningGroups = sets.New("", cmdgroups.Config)
	skipStateHandlingGroups     = sets.New("", cmdgroups.Config, cmdgroups.Meta)
)

func NewKubeswitcherCommand(opts ...KubeswitcherCommandOption) *cobra.Command {
	options := &KubeswitcherCommandOptions{}
	for _, opt := range opts {
		opt.Apply(options)
	}

	res := &cobra.Command{
		Use:   "kw <command>",
		Short: "Quickly switch between multiple Kubernetes clusters",
		Long: `Quickly switch between multiple Kubernetes clusters.

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
`,
	}

	res.AddGroup(&cobra.Group{ID: cmdgroups.Basic, Title: "Basic Commands:"})
	res.AddGroup(&cobra.Group{ID: cmdgroups.Meta, Title: "Meta Commands:"})
	res.AddGroup(&cobra.Group{ID: cmdgroups.Config, Title: "Config Commands:"})
	for _, sc := range builtinSubcommands {
		res.AddCommand(sc)
	}

	// misc
	res.SetOut(os.Stdout)
	res.SetErr(os.Stderr)
	res.SetIn(os.Stdin)

	// add plugin commands
	if !options.pluginsDisabled && len(config.Runtime.Config().Plugins) > 0 {
		res.AddGroup(&cobra.Group{ID: cmdgroups.Plugin, Title: "Plugin Commands:"})
		for _, pc := range config.Runtime.Config().Plugins {
			res.AddCommand(commandFromPluginConfig(pc))
		}
	}

	res.DisableAutoGenTag = true
	res.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// print session dir and config as debug info
		debug.Debug("Session dir: %s\n", config.Runtime.SessionDir())
		debug.Debug("Config: \n%s\n", config.Runtime.Config().String())

		// store current config in temporary history
		if err := storage.StoreToTmpHistory(); err != nil {
			libutils.Fatal(1, "error storing current config in temporary history: %w\n", err)
		}

		// rename notification message file to backup file, if it exists
		if err := fs.FS.Rename(config.Runtime.NotificationMessagePath(), config.Runtime.NotificationMessageBackupPath()); err != nil {
			if !vfs.IsNotExist(err) {
				libutils.Fatal(1, "error renaming notification message file: %w\n", err)
			}
			debug.Debug("Notification message file does not exist.\n")
		} else {
			debug.Debug("Notification message file '%s' renamed to '%s'.\n", config.Runtime.NotificationMessagePath(), config.Runtime.NotificationMessageBackupPath())
		}

		if !skipStateHandlingGroups.Has(getCmdGroup(cmd)) {
			// get last modified time of plugin state file, if it exists
			if fi, err := fs.FS.Stat(config.Runtime.PluginStatePath()); err != nil {
				if !vfs.IsNotExist(err) {
					libutils.Fatal(1, "error accessing plugin state file: %w\n", err)
				}
				debug.Debug("Plugin state file does not exist.\n")
			} else {
				pluginStateLastModified = fi.ModTime()
				debug.Debug("Plugin state file last modified: %s\n", pluginStateLastModified.Format(time.RFC3339))
			}
		}

		t, idx := internalCallStack.Peek()
		if t == nil {
			// this should only happen during the initial call
			debug.Debug("Adding initial call to internal call stack")
			t, idx = internalCallStack.Push(newTask(fmt.Sprintf("%s %s", cmd.Name(), strings.Join(args, " "))))
		}
		expectedCommand := fmt.Sprintf("%s %s", cmd.Name(), strings.Join(args, " "))
		if t.CommandArgs != expectedCommand {
			libutils.Fatal(1, "internal error: command at index %d on internal call stack doesn't match the currently executed command\nExpected: '%s'\nActual: '%s'\n", idx, expectedCommand, t.CommandArgs)
		}
		msgString := fmt.Sprintf("=== Executing call (index: %d): %s ===", idx, t.CommandArgs)
		msgSeparator := strings.Repeat("=", len(msgString))
		debug.Debug("\n%s\n%s\n%s\n", msgSeparator, msgString, msgSeparator)
	}
	res.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		t, idx := internalCallStack.Peek()
		if t == nil {
			libutils.Fatal(1, "internal error: internal call stack is empty in PersistentPostRun\n")
		}
		msgString := fmt.Sprintf("=== Finished execution of call (index: %d): %s ===", idx, t.CommandArgs)
		msgSeparator := strings.Repeat("=", len(msgString))
		debug.Debug("\n%s\n%s\n%s\n", msgSeparator, msgString, msgSeparator)
		debug.Debug("=== Starting potential internal call resulting from call at index %d ===", idx)
		if err := handleInternalCall(); err != nil {
			libutils.Fatal(1, "error handling internal call: %w\n", err)
		}
		debug.Debug("=== Finished potential internal call resulting from call at index %d ===", idx)

		cmdGroupID := getCmdGroup(cmd)

		// check if notification file exists
		// if yes, store current config in history, write generic state, delete plugin state unless it was changed,
		// copy temporary history to regular history, print notification

		noteFound := false
		_, err := fs.FS.Stat(config.Runtime.NotificationMessagePath())
		if err != nil {
			if !vfs.IsNotExist(err) {
				libutils.Fatal(1, "error accessing notification message file: %w\n", err)
			}
			debug.Debug("No notification message file found, skipping state handling.")
		} else {
			noteFound = true
			debug.Debug("Notification message file found.\n")
		}

		if noteFound {
			if skipStateHandlingGroups.Has(cmdGroupID) {
				debug.Debug("Skipping state handling because command belongs to group '%s'.\n", cmd.GroupID)
			} else {
				// write generic state
				debug.Debug("Writing generic state to '%s'.\n", config.Runtime.GenericStatePath())
				pluginName := strings.SplitN(t.CommandArgs, " ", 2)[0]
				cmdExec := fmt.Sprintf("%s %s", os.Args[0], t.CommandArgs)
				debug.Debug("\tExecuted command: %s\n", cmdExec)
				debug.Debug("\tPlugin Name: %s\n", pluginName)
				if err := state.WriteGenericState(config.Runtime.GenericStatePath(), cmdExec, pluginName); err != nil {
					libutils.Fatal(1, "error writing generic state: %w\n", err)
				}

				// delete plugin state if it was not changed
				if fi, err := fs.FS.Stat(config.Runtime.PluginStatePath()); err != nil {
					if !vfs.IsNotExist(err) {
						libutils.Fatal(1, "error accessing plugin state file: %w\n", err)
					}
					debug.Debug("Plugin state file does not exist.\n")
				} else {
					modTime := fi.ModTime()
					debug.Debug("Plugin state file last modified: %s\n", modTime.Format(time.RFC3339))
					if modTime.Equal(pluginStateLastModified) {
						debug.Debug("Plugin state file was not changed, deleting it.\n")
						if err := fs.FS.Remove(config.Runtime.PluginStatePath()); err != nil {
							libutils.Fatal(1, "error deleting plugin state file: %w\n", err)
						}
					} else {
						debug.Debug("Plugin state file was changed, keeping it.\n")
					}
				}
			}
		}

		debug.Debug("=== Starting potential callback to call at index %d ===", idx)
		if err := handleInternalCallback(t, idx); err != nil {
			libutils.Fatal(1, "error handling internal callback: %w\n", err)
		}
		debug.Debug("=== Finished potential callback to call at index %d ===", idx)
		if err := fs.FS.Remove(config.Runtime.InternalCallbackStatePath(strconv.Itoa(idx))); err == nil {
			debug.Debug("Removed used internal callback state file for index %d", idx)
		} else if !vfs.IsNotExist(err) {
			libutils.Fatal(1, "error removing used internal callback state file for index %d: %w\n", idx, err)
		}
		if !options.isInternalCallback {
			internalCallStack.Pop()
		}
		if options.isInternalCall || options.isInternalCallback {
			return
		}
		if !skipKubeconfigWarningGroups.Has(cmdGroupID) {
			// check if KUBECONFIG env var is set to the expected value and print a warning to stderr if not
			kcfg_env, ok := os.LookupEnv("KUBECONFIG")
			if !ok || kcfg_env != config.Runtime.KubeconfigPath() {
				cmd.PrintErrln("WARNING: The KUBECONFIG environment variable doesn't point to the path that is modified by kw.")
				cmd.PrintErrln("Run 'eval $(kw kubectl-env <shell>)' to fix this.")
				cmd.PrintErrln("Consider adding this command to your shell's startup script to avoid this warning in future.")
			}
		}

		debug.Debug("Removing all leftover internal callback files.\n")
		files, err := vfs.ReadDir(fs.FS, config.Runtime.SessionDir())
		if err != nil {
			libutils.Fatal(1, "error reading session directory: %w\n", err)
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), config.InternalCallbackFilePrefix) {
				debug.Debug("Removing internal callback file: %s\n", file.Name())
				if err := fs.FS.Remove(filepath.Join(config.Runtime.SessionDir(), file.Name())); err != nil {
					libutils.Fatal(1, "error removing internal callback file: %w\n", err)
				}
			}
		}

		if noteFound {
			debug.Debug("History handling")
			// copy current state to newest history entry
			if err := storage.StoreFromCurrentToHistory(); err != nil {
				libutils.Fatal(1, "error storing current state to history: %w\n", err)
			}
			// link from global history to local history
			if err := storage.StoreFromLocalToGlobalHistory(); err != nil {
				libutils.Fatal(1, "error storing local history entry to global history: %w\n", err)
			}

			// print notification, if enabled
			if config.Runtime.Config().Kubeswitcher.PrintInfoOnKubeconfigChange {
				note, err := vfs.ReadFile(fs.FS, config.Runtime.NotificationMessagePath())
				if err != nil {
					if vfs.IsNotExist(err) {
						debug.Debug("No notification message file found.\n")
						return
					}
					libutils.Fatal(1, "error accessing notification message file: %w\n", err)
				}
				if config.Runtime.Config().Kubeswitcher.PrintChangeInfoToStderr {
					cmd.PrintErrln(string(note))
				} else {
					cmd.Println(string(note))
				}
			}
		} else {
			debug.Debug("No notification message file found, skipping history handling and notification printing.")
		}
	}

	// flags
	// persistent flags have to be parsed manually, since flag parsing has to be disabled for plugin subcommands
	// so this is just for the help message
	res.PersistentFlags().BoolVar(&debug.PrintDebugStatements, "debug", debug.PrintDebugStatements, "Print debug information to stderr.")

	return res
}

var RootCmd *cobra.Command

func init() {
	cobra.EnableTraverseRunHooks = true

	// fill config variables for validation
	config.BuiltinSubcommands = sets.Set[string]{}
	config.BuiltinAliases = map[string]string{}
	for _, sc := range builtinSubcommands {
		scName := sc.Name()
		config.BuiltinSubcommands.Insert(scName)
		for _, alias := range sc.Aliases {
			config.BuiltinAliases[alias] = scName
		}
	}

	RootCmd = NewKubeswitcherCommand()
}

// // addPersistentPreRunFunctionToCommand additively adds one or more functions to the PersistentPreRun field of a command.
// // If multiple functions are added this way, they are executed in the order they have been added.
// func addPersistentPreRunFunctionToCommand(cmd *cobra.Command, funcs ...func(*cobra.Command, []string)) {
// 	if len(funcs) == 0 {
// 		return
// 	}
// 	if cmd.PersistentPreRun != nil {
// 		funcs = append([]func(*cobra.Command, []string){cmd.PersistentPreRun}, funcs...)
// 	}
// 	cmd.PersistentPreRun = func(c *cobra.Command, args []string) {
// 		for _, f := range funcs {
// 			f(c, args)
// 		}
// 	}
// }

// commandFromPluginConfig generates a cobra command from a plugin config
func commandFromPluginConfig(pc *config.PluginConfig) *cobra.Command {
	return &cobra.Command{
		Use:                pc.Name,
		Aliases:            pc.Aliases,
		DisableFlagParsing: true,
		GroupID:            cmdgroups.Plugin,
		Short:              pc.Short,
		Run: func(cmd *cobra.Command, args []string) {
			debug.Debug("executing plugin: %s %s\n", pc.Name, strings.Join(args, " "))
			config.Runtime.Context().SetPluginName(pc.Name)
			bin := exec.Command(pc.Binary, args...)
			// build command environment
			if bin.Env == nil {
				bin.Env = []string{}
			}
			bin.Env = append(bin.Env, os.Environ()...) // add current env vars
			debug.Debug("environment (in addition to parent process environment):\n")
			for k, v := range pc.Env { // add custom env vars
				debug.Debug("  %s=%s\n", k, v)
				bin.Env = append(bin.Env, fmt.Sprintf("%s=%s", k, v))
			}
			currentTaskIndexAsString := strconv.Itoa(internalCallStack.CurrentTaskIndex())
			for k, v := range config.Runtime.Context().EnvFromContext(pc.Name, pc.Config, config.Runtime.InternalCallbackRequestPath(currentTaskIndexAsString), config.Runtime.InternalCallbackStatePath(currentTaskIndexAsString)) { // add context env vars
				debug.Debug("  %s=%s\n", k, v)
				bin.Env = append(bin.Env, fmt.Sprintf("%s=%s", k, v))
			}

			// set channels
			bin.Stderr = cmd.ErrOrStderr()
			bin.Stdout = cmd.OutOrStdout()
			bin.Stdin = cmd.InOrStdin()

			// run command
			debug.Debug("starting plugin execution\n")
			if err := bin.Run(); err != nil {
				// plugin failed, try to restore state
				debug.Debug("plugin execution failed\n")
				debug.Debug("--- plugin fail ---")
				err2 := storage.LoadFromTmpHistory()
				if err2 != nil {
					err2 = fmt.Errorf("unable to restore previous state: %w", err2)
				}
				libutils.Fatal(1, "error running plugin '%s': %w\n", pc.Name, errors.Join(err, err2))
			}
			debug.Debug("finished plugin execution\n")
		},
	}
}

// getCmdGroup returns the group ID of the given command.
// If the command itself doesn't have a group ID, it checks its parents recursively.
func getCmdGroup(cmd *cobra.Command) string {
	if cmd.GroupID != "" {
		return cmd.GroupID
	}
	if cmd.HasParent() {
		return getCmdGroup(cmd.Parent())
	}
	return ""
}

type KubeswitcherCommandOptions struct {
	pluginsDisabled    bool
	isInternalCall     bool
	isInternalCallback bool
}

type KubeswitcherCommandOption interface {
	Apply(*KubeswitcherCommandOptions)
}

type DisablePlugins struct{}

var _ KubeswitcherCommandOption = DisablePlugins{}

func (DisablePlugins) Apply(opts *KubeswitcherCommandOptions) {
	opts.pluginsDisabled = true
}

type internalCallOption struct {
	internal bool
	callback bool
}

func internalCall(internal bool, callback bool) internalCallOption {
	return internalCallOption{internal: internal, callback: callback}
}

func (o internalCallOption) Apply(opts *KubeswitcherCommandOptions) {
	opts.isInternalCall = o.internal
	opts.isInternalCallback = o.callback
}
