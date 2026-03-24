package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	libcontext "github.com/Diaphteiros/kw/pluginlib/pkg/context"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
	"github.com/Diaphteiros/kw/pluginlib/pkg/state"
	libutils "github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

const (
	req_sid                           = "sid"
	req_sdir                          = "sdir"
	req_cfgDir                        = "cfgDir"
	req_cfg                           = "cfg"
	req_state                         = "state"
	KubeconfigFileName                = "kubeconfig"
	GenericStateFileName              = "state.json"
	PluginStateFileName               = "plugin_state.json"
	NotificationMessageFileName       = "message"
	NotificationMessageBackupFileName = "message.bak"
	IdFileName                        = "id"
	InternalCallFileName              = "internal_call"
	InternalCallbackFilePrefix        = "icb_"
	InternalCallbackRequestSuffix     = "_request"
	InternalCallbackStateSuffix       = "_state"
)

var (
	Runtime *KubeswitcherRuntime
)

func init() {
	Runtime = &KubeswitcherRuntime{}
	Runtime.req = libutils.NewRequirements()

	Runtime.req.Register(req_sid, Runtime.computeSessionID)
	Runtime.req.Register(req_sdir, Runtime.ensureSessionDir)
	Runtime.req.Register(req_cfgDir, Runtime.getKubeswitcherConfigDirectory)
	Runtime.req.Register(req_cfg, Runtime.loadConfig)
	Runtime.req.Register(req_state, Runtime.loadState)
}

// KubeswitcherRuntime holds the internal state during commands.
// Basically, it stores values like the session ID and directory, which have to be computed once per command call and are then reused.
type KubeswitcherRuntime struct {
	req libutils.Requirements

	sid    string
	sdir   string
	cfgDir string
	cfg    *Config
	state  *state.State
}

// SessionID() returns the current session ID
func (kr *KubeswitcherRuntime) SessionID() string {
	if err := kr.req.Require(req_sid); err != nil {
		libutils.Fatal(1, "unable to determine session ID: %w\n", err)
	}
	return kr.sid
}

// SessionDir returns the current session directory
func (kr *KubeswitcherRuntime) SessionDir() string {
	if err := kr.req.Require(req_sdir); err != nil {
		libutils.Fatal(1, "unable to determine session directory: %w\n", err)
	}
	return kr.sdir
}

// KubeconfigPath returns the path to the kubeconfig file
func (kr *KubeswitcherRuntime) KubeconfigPath() string {
	return filepath.Join(kr.SessionDir(), KubeconfigFileName)
}

// GenericStatePath returns the path to the state file
func (kr *KubeswitcherRuntime) GenericStatePath() string {
	return filepath.Join(kr.SessionDir(), GenericStateFileName)
}

// PluginStatePath returns the path to the plugin state file
func (kr *KubeswitcherRuntime) PluginStatePath() string {
	return filepath.Join(kr.SessionDir(), PluginStateFileName)
}

// NotificationMessagePath returns the path to the file containing the notification message
func (kr *KubeswitcherRuntime) NotificationMessagePath() string {
	return filepath.Join(kr.SessionDir(), NotificationMessageFileName)
}

// NotificationMessageBackupPath returns the path to the backup file containing the notification message
func (kr *KubeswitcherRuntime) NotificationMessageBackupPath() string {
	return filepath.Join(kr.SessionDir(), NotificationMessageBackupFileName)
}

// IdPath returns the path to the file containing the id
func (kr *KubeswitcherRuntime) IdPath() string {
	return filepath.Join(kr.SessionDir(), IdFileName)
}

// InternalCallPath returns the path to the file containing the internal call
func (kr *KubeswitcherRuntime) InternalCallPath() string {
	return filepath.Join(kr.SessionDir(), InternalCallFileName)
}

// InternalCallbackRequestPath returns the path where a request for an internal callback with the given ID can be placed
func (kr *KubeswitcherRuntime) InternalCallbackRequestPath(callbackID string) string {
	return filepath.Join(kr.SessionDir(), InternalCallbackFilePrefix+callbackID+InternalCallbackRequestSuffix)
}

// InternalCallbackStatePath returns the path where the state for an internal callback with the given ID can be read from
func (kr *KubeswitcherRuntime) InternalCallbackStatePath(callbackID string) string {
	return filepath.Join(kr.SessionDir(), InternalCallbackFilePrefix+callbackID+InternalCallbackStateSuffix)
}

// ConfigDirectory returns the path to the kubeswitcher config directory
func (kr *KubeswitcherRuntime) ConfigDirectory() string {
	if err := kr.req.Require(req_cfgDir); err != nil {
		libutils.Fatal(1, "unable to determine kubeswitcher config directory: %w\n", err)
	}
	return kr.cfgDir
}

// Config returns the kubeswitcher config
func (kr *KubeswitcherRuntime) Config() *Config {
	if err := kr.req.Require(req_cfg); err != nil {
		libutils.Fatal(1, "unable to load kubeswitcher config: %w\n", err)
	}
	return kr.cfg
}

// State returns the state
func (kr *KubeswitcherRuntime) State() *state.State {
	if err := kr.req.Require(req_state); err != nil {
		libutils.Fatal(1, "unable to load state: %w\n", err)
	}
	return kr.state
}

// Context returns the context.
// Creates a new context with an empty plugin name if no context has been created yet.
func (kr *KubeswitcherRuntime) Context() *libcontext.Context {
	con := libcontext.GetContext()
	if con == nil {
		con = libcontext.NewContext(
			kr.cfg.Kubeswitcher.KubectlBinary,
			kr.KubeconfigPath(),
			"",
			kr.GenericStatePath(),
			kr.PluginStatePath(),
			kr.NotificationMessagePath(),
			kr.IdPath(),
			kr.InternalCallPath(),
			"",
			"",
			"",
			kr.SessionID(),
			kr.SessionDir(),
			kr.ConfigDirectory(),
		)
	}
	return con
}

func (kr *KubeswitcherRuntime) computeSessionID() error {
	if value, ok := os.LookupEnv(ENV_KW_SESSION_ID); ok {
		if SIDRegex.MatchString(value) {
			kr.sid = value
			return nil
		}

		return fmt.Errorf("environment variable %s must only contain alphanumeric characters, underscore and dash and have a minimum length of 1 and a maximum length of 128", ENV_KW_SESSION_ID)
	}

	if value, ok := os.LookupEnv(ENV_DEFAULT_SESSION_ID); ok {
		match := UUIDRegex.FindStringSubmatch(strings.ToLower(value))
		if len(match) > 1 {
			kr.sid = match[1]
			return nil
		}
	}

	return fmt.Errorf("environment variable %s is required. Please run the following command to export it:\n\nexport %s=$(uuidgen)", ENV_KW_SESSION_ID, ENV_KW_SESSION_ID)
}

func (kr *KubeswitcherRuntime) ensureSessionDir() error {
	kr.sdir = filepath.Join(os.TempDir(), "kubeswitcher", kr.SessionID()) // don't use vfs' TempDir, as it is not stable across multiple calls
	return fs.FS.MkdirAll(kr.sdir, os.ModePerm|os.ModeDir)
}

func (kr *KubeswitcherRuntime) getKubeswitcherConfigDirectory() error {
	if value, ok := os.LookupEnv(ENV_KW_CONFIG_REPO); ok {
		kr.cfgDir = value
	} else {
		kr.cfgDir = filepath.Join(os.Getenv("HOME"), KW_CONFIG_REPO_DEFAULT_NAME)
	}

	return fs.FS.MkdirAll(kr.cfgDir, os.ModePerm|os.ModeDir)
}

func (kr *KubeswitcherRuntime) loadConfig() error {
	debug.Debug("Loading config.\n")
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	kr.cfg = cfg
	if kr.cfg == nil {
		kr.cfg = &Config{}
	}
	return nil
}

func (kr *KubeswitcherRuntime) loadState() error {
	debug.Debug("Loading state.\n")
	s, err := state.LoadState(kr.GenericStatePath(), kr.PluginStatePath())
	if err != nil {
		return err
	}
	kr.state = s
	if kr.state == nil {
		kr.state = &state.State{}
	}
	return nil
}
