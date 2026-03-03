package errors

import (
	"errors"
	"fmt"
)

type StateFromAnotherPluginError struct {
	Expected string
	Actual   string
}

func (e *StateFromAnotherPluginError) Error() string {
	return fmt.Sprintf("expected state from plugin '%s', but state is actually from plugin '%s'", e.Expected, e.Actual)
}

// NewStateFromAnotherPluginError can be use to create an error indicating that the plugin state is from another plugin than expected.
func NewStateFromAnotherPluginError(expected string, actual string) error {
	return &StateFromAnotherPluginError{
		Expected: expected,
		Actual:   actual,
	}
}

// IsStateFromAnotherPluginError returns true if the given error either is or wraps a StateFromAnotherPluginError.
func IsStateFromAnotherPluginError(err error) bool {
	_, ok := ToStateFromAnotherPluginError(err)
	return ok
}

// ToStateFromAnotherPluginError tries to parse the given error or any error wrapped by it as a StateFromAnotherPluginError.
// Returns the first StateFromAnotherPluginError found and true if such an error is found, nil and false otherwise.
func ToStateFromAnotherPluginError(err error) (*StateFromAnotherPluginError, bool) {
	terr, ok := err.(*StateFromAnotherPluginError)
	if ok {
		return terr, true
	}
	// try to unwrap (Unwrap() error)
	err = errors.Unwrap(err)
	if err != nil {
		return ToStateFromAnotherPluginError(err)
	}
	// try to unwrap (Unwrap() []error)
	u, ok := err.(interface {
		Unwrap() []error
	})
	if ok {
		for _, e := range u.Unwrap() {
			if terr, ok := e.(*StateFromAnotherPluginError); ok {
				return terr, true
			}
		}
	}
	// nope
	return nil, false
}

// IgnoreStateFromAnotherPluginError returns nil if the given error is a StateFromAnotherPluginError
// and the error itself otherwise.
func IgnoreStateFromAnotherPluginError(err error) error {
	if IsStateFromAnotherPluginError(err) {
		return nil
	}
	return err
}
