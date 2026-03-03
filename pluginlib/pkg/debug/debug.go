package debug

import (
	"fmt"
	"os"
	"strings"
)

var PrintDebugStatements bool

// Debug prints the given message (prefixed with 'DEBUG: ') to stderr if the global debug flag is set.
// Appends a newline if not already there.
func Debug(format string, values ...any) {
	if PrintDebugStatements {
		sb := strings.Builder{}
		sb.WriteString("\033[37m") // gray color
		sb.WriteString("DEBUG: ")
		sb.WriteString(fmt.Sprintf(format, values...))
		sb.WriteString("\033[0m") // reset color
		str := sb.String()
		if !strings.HasSuffix(sb.String(), "\n") {
			str += "\n"
		}
		fmt.Fprint(os.Stderr, str)
	}
}
