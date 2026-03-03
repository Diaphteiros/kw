package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"

	"sigs.k8s.io/yaml"
)

type DefaultableAndValidatable interface {
	// Default sets the default values for the config.
	Default() error
	// Validate validates the config.
	Validate() error
}

// LoadConfig takes a filepath relative to the config directory (usually just a filename) and tries to load the config from there.
// If the file doesn't exist, it takes the given defaults instead and also writes them to the config directory.
// After loading, the config is defaulted and then validated.
// Arguments:
// - cfg: The config struct to load the config into. Must be a pointer (will be passed into yaml.Unmarshal).
// - configDir: The directory where the config file is located.
// - configFileName: The name of the config file.
// - defaultConfig: The default config as []byte.
func LoadConfig[T DefaultableAndValidatable](cfg T, configDir, configFileName string, defaultConfig []byte) error {
	cfgPath := filepath.Join(configDir, configFileName)
	debug.Debug("Loading config from '%s'.\n", cfgPath)
	data, err := vfs.ReadFile(fs.FS, cfgPath)
	if err != nil {
		if !vfs.IsNotExist(err) {
			return fmt.Errorf("unable to read config file '%s': %w", cfgPath, err)
		}
		if defaultConfig != nil {
			debug.Debug("No config file '%s' found, using default config.\n", configFileName)
			data = defaultConfig
			if err := vfs.WriteFile(fs.FS, cfgPath, data, os.ModePerm); err != nil {
				return fmt.Errorf("unable to write default config file '%s': %w", cfgPath, err)
			}
		} else {
			debug.Debug("No default config file given, defaulting from Default() method.\n")
		}
	}
	if data != nil {
		err = yaml.Unmarshal(data, cfg)
		if err != nil {
			return fmt.Errorf("unable to unmarshal config '%s': %w", cfgPath, err)
		}
	}
	if err := cfg.Default(); err != nil {
		return fmt.Errorf("error defaulting config '%s': %w", cfgPath, err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config '%s': %w", cfgPath, err)
	}
	return nil
}
