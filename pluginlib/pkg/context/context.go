package context

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	"sigs.k8s.io/yaml"
)

var con *Context

type Context struct {
	// path to the kubectl binary
	// defaults to "kubectl" if not set
	KubectlBinary string `json:"kubectlBinary"`
	// path to the kubeconfig
	KubeconfigPath string `json:"kubeconfigPath"`
	// name of the currently executed plugin
	CurrentPluginName string `json:"currentPluginName"`
	// path to the generic state file
	GenericStatePath string `json:"genericStatePath"`
	// path to the plugin state file
	PluginStatePath string `json:"pluginStatePath"`
	// path to the notification message file
	NotificationMessagePath string `json:"notificationMessagePath"`
	// path to the id file
	IdPath string `json:"idPath"`
	// path to the internal call file
	InternalCallPath string `json:"internalCallPath"`
	// path to the internal callback file
	InternalCallbackPath string `json:"internalCallbackPath"`
	// statically defined plugin configuration (or empty)
	PluginConfig string `json:"pluginConfig"`
	// current session id
	SessionID string `json:"sessionID"`
	// current session directory
	SessionConfigDir string `json:"sessionConfigDir"`
	// path to the kubeswitcher config directory
	ConfigDir string `json:"configDir"`
}

// GetContext returns the current context.
// Returns nil if no context has been created yet.
func GetContext() *Context {
	return con
}

// NewContext creates a new context object from the given parameters.
// Plugins will more likely want to call NewContextFromEnv() instead.
// Note that this overwrites the current context and should only be called if GetContext() returns nil.
func NewContext(kubectlBinary, kubeconfigPath, currentPluginName, genericStatePath, pluginStatePath, notificationMessagePath, idPath, internalCallPath, internalCallbackPath, pluginConfig, sessionID, sessionConfigDir, configDir string) *Context {
	con = &Context{
		KubectlBinary:           kubectlBinary,
		KubeconfigPath:          kubeconfigPath,
		CurrentPluginName:       currentPluginName,
		GenericStatePath:        genericStatePath,
		PluginStatePath:         pluginStatePath,
		NotificationMessagePath: notificationMessagePath,
		IdPath:                  idPath,
		InternalCallPath:        internalCallPath,
		InternalCallbackPath:    internalCallbackPath,
		PluginConfig:            pluginConfig,
		SessionID:               sessionID,
		SessionConfigDir:        sessionConfigDir,
		ConfigDir:               configDir,
	}
	return con
}

// NewContextFromEnv creates a new context object from the environment variables.
// Returns an error if any of the required environment variables are not set.
// Note that this overwrites the current context and should only be called if GetContext() returns nil.
// Also sets debug.PrintDebugStatements to true if the environment variable DEBUG is set to "true".
func NewContextFromEnv() (*Context, error) {
	con = &Context{}
	ok := false
	missingEnvVarError := "env var '%s' not set"
	con.KubectlBinary, ok = os.LookupEnv(ENV_VAR_KUBECTL_PATH)
	if !ok {
		con.KubectlBinary = "kubectl"
	}
	con.KubeconfigPath, ok = os.LookupEnv(ENV_VAR_KUBECONFIG_PATH)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_KUBECONFIG_PATH)
	}
	con.CurrentPluginName, ok = os.LookupEnv(ENV_VAR_CURRENT_PLUGIN_NAME)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_CURRENT_PLUGIN_NAME)
	}
	con.GenericStatePath, ok = os.LookupEnv(ENV_VAR_GENERIC_STATE_PATH)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_GENERIC_STATE_PATH)
	}
	con.PluginStatePath, ok = os.LookupEnv(ENV_VAR_PLUGIN_STATE_PATH)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_PLUGIN_STATE_PATH)
	}
	con.NotificationMessagePath, ok = os.LookupEnv(ENV_VAR_NOTIFICATION_MESSAGE_PATH)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_NOTIFICATION_MESSAGE_PATH)
	}
	con.IdPath, ok = os.LookupEnv(ENV_VAR_ID_PATH)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_ID_PATH)
	}
	con.InternalCallPath, ok = os.LookupEnv(ENV_VAR_INTERNAL_CALL_PATH)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_INTERNAL_CALL_PATH)
	}
	con.InternalCallbackPath, ok = os.LookupEnv(ENV_VAR_INTERNAL_CALLBACK_PATH)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_INTERNAL_CALL_PATH)
	}
	con.PluginConfig = os.Getenv(ENV_VAR_PLUGIN_CONFIG) // plugin config is optional
	con.SessionID, ok = os.LookupEnv(ENV_VAR_SESSION_ID)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_SESSION_ID)
	}
	con.SessionConfigDir, ok = os.LookupEnv(ENV_VAR_SESSION_CONFIG_DIR)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_SESSION_CONFIG_DIR)
	}
	con.ConfigDir, ok = os.LookupEnv(ENV_VAR_CONFIG_DIR)
	if !ok {
		return nil, fmt.Errorf(missingEnvVarError, ENV_VAR_CONFIG_DIR)
	}
	rawDebug, ok := os.LookupEnv(ENV_VAR_DEBUG)
	if ok && rawDebug == "true" {
		debug.PrintDebugStatements = true
	}
	return con, nil
}

// EnvFromContext generates a map of environment variables from the context.
// Plugin name and config are injected, since they are usually not set when this is called.
func (con *Context) EnvFromContext(pluginName string, pluginConfig []byte, internalCallbackPath string) map[string]string {
	env := map[string]string{
		ENV_VAR_KUBECTL_PATH:              con.KubectlBinary,
		ENV_VAR_KUBECONFIG_PATH:           con.KubeconfigPath,
		ENV_VAR_CURRENT_PLUGIN_NAME:       pluginName,
		ENV_VAR_GENERIC_STATE_PATH:        con.GenericStatePath,
		ENV_VAR_PLUGIN_STATE_PATH:         con.PluginStatePath,
		ENV_VAR_NOTIFICATION_MESSAGE_PATH: con.NotificationMessagePath,
		ENV_VAR_ID_PATH:                   con.IdPath,
		ENV_VAR_INTERNAL_CALL_PATH:        con.InternalCallPath,
		ENV_VAR_INTERNAL_CALLBACK_PATH:    internalCallbackPath,
		ENV_VAR_SESSION_ID:                con.SessionID,
		ENV_VAR_SESSION_CONFIG_DIR:        con.SessionConfigDir,
		ENV_VAR_CONFIG_DIR:                con.ConfigDir,
		ENV_VAR_DEBUG:                     fmt.Sprintf("%t", os.Getenv(ENV_VAR_DEBUG) == "true" || debug.PrintDebugStatements),
	}
	if pluginConfig != nil {
		env[ENV_VAR_PLUGIN_CONFIG] = string(pluginConfig)
	}
	return env
}

func (con *Context) String() string {
	data, err := yaml.Marshal(con)
	if err != nil {
		return fmt.Sprintf("unable to marshal context: %v", err)
	}
	return string(data)
}

// SetPluginName sets the name of the currently executed plugin, if it is empty.
// No-op if the name is already set.
func (con *Context) SetPluginName(name string) {
	if con.CurrentPluginName == "" {
		con.CurrentPluginName = name
	}
}

// Writes the kubeconfig and a notification message to the respective files.
// This assumes that the kubeconfig has actually changed.
// The additional args are passed into fmt.Sprintf with the message for formatting.
func (con *Context) WriteKubeconfig(kcfg []byte, message string, args ...any) error {
	// write the kubeconfig
	debug.Debug("Writing kubeconfig to %s\n", con.KubeconfigPath)
	// removing file before in case it is a symlink
	if err := fs.FS.Remove(con.KubeconfigPath); err != nil && !vfs.IsNotExist(err) {
		return fmt.Errorf("unable to remove kubeconfig: %w", err)
	}
	if err := vfs.WriteFile(fs.FS, con.KubeconfigPath, kcfg, os.ModePerm); err != nil {
		return fmt.Errorf("unable to write kubeconfig: %w", err)
	}

	// write the notification message
	return con.WriteNotificationMessage(message, args...)
}

func (con *Context) WriteKubeconfigSymlink(kcfgPath, message string, args ...any) error {
	// write the kubeconfig
	debug.Debug("Writing kubeconfig to %s\n", con.KubeconfigPath)
	// removing file before in case it is a symlink
	if err := fs.FS.Remove(con.KubeconfigPath); err != nil && !vfs.IsNotExist(err) {
		return fmt.Errorf("unable to remove kubeconfig: %w", err)
	}
	if err := fs.FS.Symlink(kcfgPath, con.KubeconfigPath); err != nil {
		return fmt.Errorf("unable to write kubeconfig: %w", err)
	}

	// write the notification message
	return con.WriteNotificationMessage(message, args...)
}

// WriteNotificationMessage writes the given message to the notification message file.
// The additional args are passed into fmt.Sprintf with the message for formatting.
// Note that this function is usually called by WriteKubeconfig (or WriteKubeconfigSymlink) and does not need to be called manually.
func (con *Context) WriteNotificationMessage(message string, args ...any) error {
	debug.Debug("Writing notification message to %s\n", con.NotificationMessagePath)
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	if err := vfs.WriteFile(fs.FS, con.NotificationMessagePath, []byte(message), os.ModePerm); err != nil {
		return fmt.Errorf("unable to write notification message: %w", err)
	}

	return nil
}

// WritePluginState writes the given state to the plugin state file.
// If the parameter has the type []byte, it is written as is.
// Otherwise, it is marshaled into JSON before being written.
// If this function has not been called on the context when Close() is called,
// the plugin state will be deleted to ensure that no leftover state from a previous command is kept.
func (con *Context) WritePluginState(ps any) error {
	// write the plugin state
	// check if state is given as a byte slice
	if psb, ok := ps.([]byte); ok {
		debug.Debug("Writing byte-slice plugin state to %s\n", con.PluginStatePath)
		err := vfs.WriteFile(fs.FS, con.PluginStatePath, psb, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to write plugin state: %w", err)
		}
		return nil
	}

	// otherwise, marshal into json
	debug.Debug("Writing object plugin state to %s\n", con.PluginStatePath)
	psj, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal plugin state into json: %w", err)
	}
	err = vfs.WriteFile(fs.FS, con.PluginStatePath, psj, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to write plugin state: %w", err)
	}
	return nil
}

// WriteId formats the given string and writes it to the id file.
// All newlines are removed from the string.
func (con *Context) WriteId(id string, args ...any) error {
	// write the id
	debug.Debug("Writing id to %s\n", con.IdPath)
	if len(args) > 0 {
		id = fmt.Sprintf(id, args...)
	}
	err := vfs.WriteFile(fs.FS, con.IdPath, []byte(strings.ReplaceAll(id, "\n", "")), os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to write id: %w", err)
	}

	return nil
}

// WriteInternalCall requests an internal call by writing the command to the internal call file.
// Note that the name of the kubeswitcher binary is omitted, so to call 'kw foo --bar', the call parameter should be 'foo --bar'.
// If callback is not nil, the calling command will receive a callback after the internal call has been executed.
// During the callback, the information written here can be read and used to react to the internal call.
func (con *Context) WriteInternalCall(call string, callback []byte) error {
	debug.Debug("Writing internal call to %s\n", con.InternalCallPath)
	err := vfs.WriteFile(fs.FS, con.InternalCallPath, []byte(call), os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to write internal call: %w", err)
	}

	if callback == nil {
		debug.Debug("No internal callback to write\n")
	} else {
		debug.Debug("Writing internal callback to %s\n", con.InternalCallbackPath)
		err = vfs.WriteFile(fs.FS, con.InternalCallbackPath, callback, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to write internal callback: %w", err)
		}
	}

	return nil
}
