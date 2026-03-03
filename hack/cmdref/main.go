package main

import (
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/Diaphteiros/kw/cmd"
)

func main() {
	if len(os.Args) < 2 {
		panic("documentation folder path required as argument")
	}
	// disable plugins to avoid having plugin-specific subcommands in the general command reference
	if err := doc.GenMarkdownTree(cmd.NewKubeswitcherCommand(cmd.DisablePlugins{}), os.Args[1]); err != nil {
		panic(err)
	}
}
