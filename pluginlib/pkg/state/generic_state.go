package state

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
)

type GenericState struct {
	// LastUsed holds information about the last used kubeswitcher command
	LastUsed *LastUsed `json:"lastUsed,omitempty"`
}

// LastUsed holds information about the last used kubeswitcher command
type LastUsed struct {
	// Command is the last used command
	Command string `json:"command,omitempty"`
	// Plugin is the name of the plugin that handled the last command
	Plugin string `json:"plugin,omitempty"`
}

// WriteGenericState builds a GenericState object from the given parameters and writes it to the given path as json.
func WriteGenericState(path string, lastCommand string, lastPlugin string) error {
	state := &GenericState{
		LastUsed: &LastUsed{
			Command: lastCommand,
			Plugin:  lastPlugin,
		},
	}

	data, err := json.MarshalIndent(state, "", "	")
	if err != nil {
		return fmt.Errorf("unable to marshal state into json: %w", err)
	}
	if err := vfs.WriteFile(fs.FS, path, data, os.ModePerm); err != nil {
		return fmt.Errorf("unable to write state: %w", err)
	}
	return nil
}
