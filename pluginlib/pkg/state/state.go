package state

import (
	"fmt"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"sigs.k8s.io/yaml"

	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
)

// State holds information about the current state.
// This is stored in files between command calls and contains information like the last used command.
type State struct {
	GenericState `json:"genericState"`
	// RawPluginState holds the plugin-specific state.
	// As this needs to be parsed by the plugin itself, the file is just read and dumped into this field.
	RawPluginState []byte `json:"pluginState"`

	// metadata, mostly for reloading the state
	genericStatePath string
	pluginStatePath  string
}

// LoadState reads the state from the given files.
func LoadState(genericStatePath, pluginStatePath string) (*State, error) {
	debug.Debug("Loading state from '%s' and '%s'.\n", genericStatePath, pluginStatePath)
	s := &State{}
	if err := loadState_helper(s, genericStatePath, pluginStatePath); err != nil {
		return nil, fmt.Errorf("unable to load state: %w", err)
	}
	if s.genericStatePath == "" {
		return nil, nil // no state yet
	}
	return s, nil
}

// loadState_helper is an auxiliary function that reads the state into the given state object.
// This is to unify the coding for LoadState and Reload.
func loadState_helper(s *State, genericStatePath, pluginStatePath string) error {
	// generic state
	data, err := vfs.ReadFile(fs.FS, genericStatePath)
	if err != nil {
		if !vfs.IsNotExist(err) {
			return fmt.Errorf("unable to read state file '%s': %w", genericStatePath, err)
		}
		return nil // no state yet
	}
	if err := yaml.Unmarshal(data, &s.GenericState); err != nil {
		return fmt.Errorf("unable to unmarshal generic state: %w", err)
	}
	s.genericStatePath = genericStatePath

	// plugin state
	s.RawPluginState, err = vfs.ReadFile(fs.FS, pluginStatePath)
	if err != nil {
		if !vfs.IsNotExist(err) {
			return fmt.Errorf("unable to read plugin state file '%s': %w", pluginStatePath, err)
		}
		return nil // no plugin state
	}
	s.pluginStatePath = pluginStatePath

	return nil
}

// Reload updates the state object with the current state from the files.
func (s *State) Reload() error {
	return loadState_helper(s, s.genericStatePath, s.pluginStatePath)
}
