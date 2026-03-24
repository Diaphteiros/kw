package debug

import (
	"fmt"
	"os"
	"strings"
)

var PrintDebugStatements bool

// Debug prints the given message to stderr if the global debug flag is set.
// Appends a newline.
func Debug(format string, values ...any) {
	if PrintDebugStatements {
		sb := strings.Builder{}
		sb.WriteString("\033[37m") // gray color
		fmt.Fprintf(&sb, format, values...)
		sb.WriteString("\n")
		sb.WriteString("\033[0m") // reset color
		str := sb.String()
		fmt.Fprint(os.Stderr, str)
	}
}
