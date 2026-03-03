package state

import (
	"fmt"

	"sigs.k8s.io/yaml"

	"github.com/Diaphteiros/kw/pluginlib/pkg/errors"
)

type TypedState[T any] struct {
	*State
	PluginState T
}

// LoadTypedState loads the state and then attempts to unmarshal the plugin state into the given type.
// If expectedPluginName is non-empty, the function will throw an StateFromAnotherPluginError if the last used plugin differs from the expected one.
func LoadTypedState[T any](genericStatePath, pluginStatePath, expectedPluginName string) (*TypedState[T], error) {
	ts := &TypedState[T]{}
	var err error
	ts.State, err = LoadState(genericStatePath, pluginStatePath)
	if err != nil {
		return nil, err
	}
	if ts.State == nil {
		return nil, nil
	}
	if expectedPluginName != "" && (ts.LastUsed == nil || ts.LastUsed.Plugin != expectedPluginName) {
		return nil, errors.NewStateFromAnotherPluginError(expectedPluginName, ts.LastUsed.Plugin)
	}
	if len(ts.RawPluginState) > 0 {
		err = yaml.Unmarshal(ts.RawPluginState, &ts.PluginState)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling plugin state: %w", err)
		}
	}
	return ts, nil
}
