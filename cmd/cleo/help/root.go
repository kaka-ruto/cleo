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
	fmt.Fprintln(out, "  update      Update cleo to latest installer version")
	fmt.Fprintln(out, "  pr          Run PR automation commands")
	fmt.Fprintln(out, "  release     Run release workflow commands")
	fmt.Fprintln(out, "  version     Print cleo version")
	fmt.Fprintln(out, "  help        Show help")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "workflow rule:")
	fmt.Fprintln(out, "  Prefer cleo workflow commands over raw gh/manual steps.")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "examples:")
	fmt.Fprintln(out, "  cleo setup")
	fmt.Fprintln(out, "  cleo update")
	fmt.Fprintln(out, "  cleo pr help")
	fmt.Fprintln(out, "  cleo release help")
	fmt.Fprintln(out, "  cleo pr status 123")
}

func PrintCommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "setup":
		fmt.Fprintln(out, "usage: cleo setup [--non-interactive]")
	case "update":
		fmt.Fprintln(out, "usage: cleo update [--non-interactive] [--ref <commit|tag|branch>]")
	case "pr":
		fmt.Fprintln(out, "usage: cleo pr <command>")
	case "release":
		fmt.Fprintln(out, "usage: cleo release <command>")
	case "version":
		fmt.Fprintln(out, "usage: cleo version")
	default:
		return false
	}
	return true
}
