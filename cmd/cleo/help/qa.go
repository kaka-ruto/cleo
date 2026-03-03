package help

import (
	"fmt"
	"io"
)

func PrintQA(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo qa <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  init                                     Install/update reusable QA kit files")
	fmt.Fprintln(out, "  scaffold [--title <text>]              Print BDD AC scaffold YAML")
	fmt.Fprintln(out, "  start --source <branch|pr|request> --ref <name|id|text> --goals <text> [--ac <yaml>]")
	fmt.Fprintln(out, "  doctor --session <id>")
	fmt.Fprintln(out, "  plan --session <id>")
	fmt.Fprintln(out, "  run --session <id> [--mode <auto|manual|pr>]")
	fmt.Fprintln(out, "  log --session <id> --title <text> --details <text> [--severity <low|medium|high|critical>]")
	fmt.Fprintln(out, "  finish --session <id> --verdict <pass|fail|blocked>")
	fmt.Fprintln(out, "  report --session <id> [--publish <pr>] [--ref <pr>]")
	fmt.Fprintln(out, "  help [command]")
}

func PrintQACommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "init":
		fmt.Fprintln(out, "usage: cleo qa init")
	case "scaffold":
		fmt.Fprintln(out, "usage: cleo qa scaffold [--title <text>]")
	case "start":
		fmt.Fprintln(out, "usage: cleo qa start --source <branch|pr|request> --ref <name|id|text> --goals <text> [--ac <yaml>]")
	case "doctor":
		fmt.Fprintln(out, "usage: cleo qa doctor --session <id>")
	case "plan":
		fmt.Fprintln(out, "usage: cleo qa plan --session <id>")
	case "run":
		fmt.Fprintln(out, "usage: cleo qa run --session <id> [--mode <auto|manual|pr>]")
	case "log":
		fmt.Fprintln(out, "usage: cleo qa log --session <id> --title <text> --details <text> [--severity <low|medium|high|critical>]")
	case "finish":
		fmt.Fprintln(out, "usage: cleo qa finish --session <id> --verdict <pass|fail|blocked>")
	case "report":
		fmt.Fprintln(out, "usage: cleo qa report --session <id> [--publish <pr>] [--ref <pr>]")
	default:
		return false
	}
	return true
}
