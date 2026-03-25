package selector

import (
	"errors"
	"fmt"
	"slices"

	fuzzy "github.com/ktr0731/go-fuzzyfinder"
	"sigs.k8s.io/yaml"

	"github.com/Diaphteiros/kw/pluginlib/pkg/utils"
)

type Selector[T any] struct {
	fuzzyArgs           []fuzzy.Option
	getElem             func(idx int) T
	keyFunc             func(elem T) string
	data                []T
	sortFunc            func(a, b T) int
	fatalOnAbortMessage string
	fatalOnErrorMessage string
}

func New[T any]() *Selector[T] {
	return &Selector[T]{
		fuzzyArgs: []fuzzy.Option{},
	}
}

// WithFuzzyArgs adds arbitrary fuzzy finder options to the selector.
// It is recommended to only use this option if the other functions do not cover the desired functionality.
func (s *Selector[T]) WithFuzzyArgs(args ...fuzzy.Option) *Selector[T] {
	s.fuzzyArgs = append(s.fuzzyArgs, args...)
	return s
}

// WithPrompt sets the prompt string for the fuzzy finder.
func (s *Selector[T]) WithPrompt(prompt string) *Selector[T] {
	s.fuzzyArgs = append(s.fuzzyArgs, fuzzy.WithPromptString(prompt))
	return s
}

// WithQuery sets the initial query string for the fuzzy finder.
func (s *Selector[T]) WithQuery(query string) *Selector[T] {
	s.fuzzyArgs = append(s.fuzzyArgs, fuzzy.WithQuery(query))
	return s
}

// WithPreview sets the preview window function for the fuzzy finder.
// The preview function is called with the selected element and the width and height of the preview window.
func (s *Selector[T]) WithPreview(previewFunc func(elem T, width, height int) string) *Selector[T] {
	s.fuzzyArgs = append(s.fuzzyArgs, fuzzy.WithPreviewWindow(func(i, width, height int) string {
		if i < 0 || i >= len(s.data) {
			return "No entry selected"
		}
		elem := s.getElem(i)
		return previewFunc(elem, width, height)
	}))
	return s
}

// WithYamlPreview is an alternative to WithPreview and simply renders the complete selected item as YAML in the preview window.
func (s *Selector[T]) WithYamlPreview() *Selector[T] {
	s.fuzzyArgs = append(s.fuzzyArgs, fuzzy.WithPreviewWindow(func(i, width, height int) string {
		if i < 0 || i >= len(s.data) {
			return "No entry selected"
		}
		if s.getElem == nil {
			return "internal error: no getElem function set"
		}
		item, err := yaml.Marshal(s.getElem(i))
		if err != nil {
			return fmt.Sprintf("Error rendering entry: %v", err)
		}
		return string(item)
	}))
	return s
}

// WithSortFunc sets a sorting function which is applied to the data before displaying it in the fuzzy finder.
// Works like the standard sorting logic, the function should return a negative number if a < b, a positive number if a > b and 0 if they are equal.
func (s *Selector[T]) WithSortFunc(sortFunc func(a, b T) int) *Selector[T] {
	s.sortFunc = sortFunc
	return s
}

// WithFatalOnAbort sets the message that is displayed when the user aborts the selection or does not select a valid entry.
// If this is not empty, the Select method will fatal with this message when the user aborts the selection.
func (s *Selector[T]) WithFatalOnAbort(msg string) *Selector[T] {
	s.fatalOnAbortMessage = msg
	return s
}

// WithFatalOnError sets the message that is displayed when an error occurs during selection.
// If this is not empty, the Select method will fatal with this message when an error occurs.
// The message is expected to be a format string with a single %w verb for the error.
func (s *Selector[T]) WithFatalOnError(msg string) *Selector[T] {
	s.fatalOnErrorMessage = msg
	return s
}

// From sets the slice from which to select and the function that computes the key (which is displayed) of a given item.
// The keyFunc should return a string that represents the passed in entry.
// Duplicate keys are allowed, but may be indistinguishable if not further distinguished by a preview function.
func (s *Selector[T]) From(data []T, keyFunc func(elem T) string) *Selector[T] {
	s.data = make([]T, len(data))
	copy(s.data, data)
	if s.sortFunc != nil {
		slices.SortStableFunc(s.data, s.sortFunc)
	}
	s.keyFunc = keyFunc
	s.getElem = func(idx int) T {
		if idx < 0 || idx >= len(s.data) {
			var zero T
			return zero
		}
		return s.data[idx]
	}
	return s
}

// Select launches the fuzzy finder for the actual selection.
// It returns the index and the selected item.
// Note that the function will not return an error and exit instead if the 'WithFatalOn...' method was called for the respective error case.
func (s *Selector[T]) Select() (int, T, error) {
	// the user either did not enter an key or it did not match any existing key
	// use the fuzzy finder to select one
	idx, err := fuzzy.Find(s.data, func(i int) string {
		return s.keyFunc(s.getElem(i))
	},
		s.fuzzyArgs...,
	)
	if err != nil {
		if errors.Is(err, fuzzy.ErrAbort) {
			if s.fatalOnAbortMessage != "" {
				utils.Fatal(1, s.fatalOnAbortMessage)
			}
		} else {
			if s.fatalOnErrorMessage != "" {
				utils.Fatal(1, s.fatalOnErrorMessage, err)
			}
		}
	}
	return idx, s.getElem(idx), err
}

// Identity can be used as a key function argument to the selector's 'From' method.
// It requires the generic type to be a string-like (which can be cast to a string) and simply returns the string itself as the key.
func Identity[T ~string](raw T) string {
	return string(raw)
}

// Invert wraps a sorting function and inverts its result, thereby reversing the sorting order.
func Invert[T any](sortFunc func(a, b T) int) func(a, b T) int {
	return func(a, b T) int {
		return -sortFunc(a, b)
	}
}
