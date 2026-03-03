package help

import (
	"fmt"
	"io"
)

func PrintPR(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo pr <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "workflow rule: use cleo pr commands for PR work.")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  status <pr>                            Show PR summary")
	fmt.Fprintln(out, "  gate <pr>                              Validate PR is merge-ready")
	fmt.Fprintln(out, "  checks <pr>                            Show GitHub check status")
	fmt.Fprintln(out, "  watch <pr|sha>                         Watch checks until completion")
	fmt.Fprintln(out, "  doctor                                 Validate local PR tooling")
	fmt.Fprintln(out, "  run <pr> [--dry]                       Run post-merge commands")
	fmt.Fprintln(out, "  create [flags]                         Create PR with cleo template")
	fmt.Fprintln(out, "  merge <pr> [flags]                     Merge PR with safety steps")
	fmt.Fprintln(out, "  rebase <pr>                            Rebase PR branch")
	fmt.Fprintln(out, "  retarget <pr> --base <branch>          Change PR base branch")
	fmt.Fprintln(out, "  batch [--from <pr>] [flags]            Merge open PRs in sequence")
	fmt.Fprintln(out, "  help [command]                         Show PR help")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "examples:")
	fmt.Fprintln(out, "  cleo pr help merge")
	fmt.Fprintln(out, "  cleo pr merge 123 --delete-branch")
}

func PrintPRCommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "status":
		fmt.Fprintln(out, "usage: cleo pr status <pr>")
	case "gate":
		fmt.Fprintln(out, "usage: cleo pr gate <pr>")
	case "checks":
		fmt.Fprintln(out, "usage: cleo pr checks <pr>")
	case "watch":
		fmt.Fprintln(out, "usage: cleo pr watch <pr|sha>")
	case "doctor":
		fmt.Fprintln(out, "usage: cleo pr doctor")
	case "run":
		fmt.Fprintln(out, "usage: cleo pr run <pr> [--dry]")
	case "create":
		fmt.Fprintln(out, "usage: cleo pr create [--title ...] [--summary ...] [--why ...] [--what ...] [--test ...] [--risk ...] [--rollback ...] [--owner ...] [--ac ...] [--cmd ...] [--draft]")
	case "merge":
		fmt.Fprintln(out, "usage: cleo pr merge <pr> [--no-watch] [--no-run] [--no-rebase] [--delete-branch]")
	case "batch":
		fmt.Fprintln(out, "usage: cleo pr batch [--from <pr>] [--no-watch] [--no-run] [--no-rebase]")
	case "rebase":
		fmt.Fprintln(out, "usage: cleo pr rebase <pr>")
	case "retarget":
		fmt.Fprintln(out, "usage: cleo pr retarget <pr> --base <branch>")
	default:
		return false
	}
	return true
}
