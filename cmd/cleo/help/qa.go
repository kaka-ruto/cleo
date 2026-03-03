package help

import (
	"fmt"
	"io"
)

func PrintQA(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo qa <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  start --source <branch|pr|request> --ref <name|id|text> --goals <text>")
	fmt.Fprintln(out, "  doctor --ac-file <path>")
	fmt.Fprintln(out, "  plan --session <id> --ac-file <path>")
	fmt.Fprintln(out, "  run --session <id> --ac-file <path>")
	fmt.Fprintln(out, "  log --session <id> --title <text> --details <text> [--severity <low|medium|high|critical>]")
	fmt.Fprintln(out, "  finish --session <id> --verdict <pass|fail|blocked>")
	fmt.Fprintln(out, "  report --session <id>")
	fmt.Fprintln(out, "  help [command]")
}

func PrintQACommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "start":
		fmt.Fprintln(out, "usage: cleo qa start --source <branch|pr|request> --ref <name|id|text> --goals <text>")
	case "doctor":
		fmt.Fprintln(out, "usage: cleo qa doctor --ac-file <path>")
	case "plan":
		fmt.Fprintln(out, "usage: cleo qa plan --session <id> --ac-file <path>")
	case "run":
		fmt.Fprintln(out, "usage: cleo qa run --session <id> --ac-file <path>")
	case "log":
		fmt.Fprintln(out, "usage: cleo qa log --session <id> --title <text> --details <text> [--severity <low|medium|high|critical>]")
	case "finish":
		fmt.Fprintln(out, "usage: cleo qa finish --session <id> --verdict <pass|fail|blocked>")
	case "report":
		fmt.Fprintln(out, "usage: cleo qa report --session <id>")
	default:
		return false
	}
	return true
}
