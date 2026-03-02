package help

import (
	"fmt"
	"io"
)

func PrintRoot(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  setup       Run setup wizard for current repository")
	fmt.Fprintln(out, "  pr          Run PR automation commands")
	fmt.Fprintln(out, "  version     Print cleo version")
	fmt.Fprintln(out, "  help        Show help")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "examples:")
	fmt.Fprintln(out, "  cleo setup")
	fmt.Fprintln(out, "  cleo pr help")
	fmt.Fprintln(out, "  cleo pr status 123")
}
