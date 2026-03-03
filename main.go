package main

import (
	"os"
	"slices"

	"github.com/Diaphteiros/kw/cmd"
	libcontext "github.com/Diaphteiros/kw/pluginlib/pkg/context"
	"github.com/Diaphteiros/kw/pluginlib/pkg/debug"
)

func main() {
	// manually parse the '--debug' flag to enable debug output
	args := os.Args[1:]
	stopIdx := slices.Index(args, "--")
	if stopIdx < 0 {
		stopIdx = len(args)
	}
	if idx := slices.Index(args[:stopIdx], "--debug"); idx >= 0 {
		args = append(args[:idx], args[idx+1:]...)
		os.Setenv(libcontext.ENV_VAR_DEBUG, "true")
		debug.PrintDebugStatements = true
	}
	cmd.RootCmd.SetArgs(args)
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
