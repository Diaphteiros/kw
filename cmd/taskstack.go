package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/Diaphteiros/kw/pkg/config"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
	"github.com/Diaphteiros/kw/pluginlib/pkg/fs"
)

var internalCallStack *taskStack = newTaskStack()

type taskStack []*task

type task struct {
	CommandArgs string
}

func newTask(commandArgs string) *task {
	return &task{
		CommandArgs: commandArgs,
	}
}

func newTaskStack() *taskStack {
	return &taskStack{}
}

// Push adds a task to the stack.
// It returns the added task and its index.
func (ts *taskStack) Push(t *task) (*task, int) {
	*ts = append(*ts, t)
	debug.Debug("Pushed '%s' to internal call stack: %s", t.CommandArgs, internalCallStack.String())
	return t, ts.Len() - 1
}

// Peek returns the last task from the stack without removing it, together with its index.
// If the stack is empty, nil and -1 are returned.
func (ts *taskStack) Peek() (*task, int) {
	idx := ts.CurrentTaskIndex()
	if idx < 0 {
		return nil, -1
	}
	return (*ts)[idx], idx
}

// Pop removes the last task from the stack and returns it together with its index.
// If the stack is empty, nil and -1 are returned.
func (ts *taskStack) Pop() (*task, int) {
	idx := len(*ts) - 1
	if idx < 0 {
		return nil, -1
	}
	t := (*ts)[idx]
	*ts = (*ts)[:idx]
	debug.Debug("Popped '%s' from internal call stack: %s", t.CommandArgs, internalCallStack.String())
	return t, idx
}

// handleInternalCall is meant to be called after the main command execution.
// It checks whether an internal call was requested, and if so, creates a task for it and executes it.
func handleInternalCall() error {
	internalCallFilePath := config.Runtime.InternalCallPath()
	internalCallRaw, err := vfs.ReadFile(fs.FS, internalCallFilePath)
	if err != nil {
		if !vfs.IsNotExist(err) {
			return fmt.Errorf("error accessing internal call file: %w", err)
		}
		debug.Debug("No internal call file found.\n")
	} else {
		if err := fs.FS.Remove(internalCallFilePath); err != nil {
			return fmt.Errorf("error deleting internal call file '%s': %w", internalCallFilePath, err)
		}
		t, _ := internalCallStack.Push(newTask(strings.TrimSpace(string(internalCallRaw))))
		if err := executeInternalCall(t.CommandArgs, false); err != nil {
			return err
		}
	}
	return nil
}

// handleInternalCallback takes a task that has just run and its index and checks whether there is a callback file for this index.
// If so, it executes the callback command.
func handleInternalCallback(t *task, idx int) error {
	callbackRequestPath := config.Runtime.InternalCallbackRequestPath(strconv.Itoa(idx))
	_, err := fs.FS.Stat(callbackRequestPath)
	if err != nil {
		if !vfs.IsNotExist(err) {
			return fmt.Errorf("error accessing internal callback file: %w", err)
		}
		debug.Debug("No internal callback file found for index %d, skipping callback execution", idx)
	} else {
		// remove any existing callback state file and rename the request file to the state file
		callbackStatePath := config.Runtime.InternalCallbackStatePath(strconv.Itoa(idx))
		if err := fs.FS.Remove(callbackStatePath); err != nil && !vfs.IsNotExist(err) {
			return fmt.Errorf("error deleting internal callback state file '%s': %w", callbackStatePath, err)
		}
		if err := fs.FS.Rename(callbackRequestPath, callbackStatePath); err != nil {
			return fmt.Errorf("error renaming internal callback request file from '%s' to '%s': %w", callbackRequestPath, callbackStatePath, err)
		}
		if err := executeInternalCall(t.CommandArgs, true); err != nil {
			return err
		}
	}
	return nil
}

func (ts *taskStack) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, t := range *ts {
		fmt.Fprintf(&sb, "'%s'", t.CommandArgs)
		if i < len(*ts)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (ts *taskStack) Len() int {
	if ts == nil {
		return 0
	}
	return len(*ts)
}

func (ts *taskStack) CurrentTaskIndex() int {
	return ts.Len() - 1
}

func executeInternalCall(commandArgs string, isCallback bool) error {
	debug.Debug("Executing internal call: %s", commandArgs)
	internalCmd := NewKubeswitcherCommand(internalCall(true, isCallback))
	internalCmd.SetArgs(strings.Split(commandArgs, " "))
	if err := internalCmd.Execute(); err != nil {
		return fmt.Errorf("error executing internal call '%s': %w", commandArgs, err)
	}
	return nil
}
