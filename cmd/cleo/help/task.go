package help

import (
	"fmt"
	"io"
)

func PrintTask(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo task <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  list [--status <open|in_progress|closed>]")
	fmt.Fprintln(out, "  show --id <task-id>")
	fmt.Fprintln(out, "  claim --id <task-id>")
	fmt.Fprintln(out, "  work --id <task-id> [--new-branch|--in-place]")
	fmt.Fprintln(out, "  close --id <task-id>")
	fmt.Fprintln(out, "  help [command]")
}

func PrintTaskCommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "list":
		fmt.Fprintln(out, "usage: cleo task list [--status <open|in_progress|closed>]")
	case "show":
		fmt.Fprintln(out, "usage: cleo task show --id <task-id>")
	case "claim":
		fmt.Fprintln(out, "usage: cleo task claim --id <task-id>")
	case "work":
		fmt.Fprintln(out, "usage: cleo task work --id <task-id> [--new-branch|--in-place]")
	case "close":
		fmt.Fprintln(out, "usage: cleo task close --id <task-id>")
	default:
		return false
	}
	return true
}
