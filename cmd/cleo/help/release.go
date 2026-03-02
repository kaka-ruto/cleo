package help

import (
	"fmt"
	"io"
)

func PrintRelease(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo release <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "workflow rule: use cleo release commands for release work.")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  list [--limit N]                       List recent releases")
	fmt.Fprintln(out, "  latest                                 Show latest release")
	fmt.Fprintln(out, "  plan --version <vX.Y.Z>                Validate release preconditions")
	fmt.Fprintln(out, "  cut --version <vX.Y.Z>                 Create and push tag")
	fmt.Fprintln(out, "  publish --version <vX.Y.Z> [flags]     Create GitHub release (Go repos auto-attach artifacts)")
	fmt.Fprintln(out, "  verify --version <vX.Y.Z>              Verify published release")
	fmt.Fprintln(out, "  go <command>                           Run explicit Go release flow")
	fmt.Fprintln(out, "  help [command]                         Show release help")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "examples:")
	fmt.Fprintln(out, "  cleo release list --limit 10")
	fmt.Fprintln(out, "  cleo release latest")
	fmt.Fprintln(out, "  cleo release plan --version v0.1.0")
	fmt.Fprintln(out, "  cleo release cut --version v0.1.0")
	fmt.Fprintln(out, "  cleo release publish --version v0.1.0 --final")
	fmt.Fprintln(out, "  cleo release go publish --version v0.1.0 --final")
}

func PrintReleaseCommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "list":
		fmt.Fprintln(out, "usage: cleo release list [--limit N]")
	case "latest":
		fmt.Fprintln(out, "usage: cleo release latest")
	case "plan":
		fmt.Fprintln(out, "usage: cleo release plan --version <vX.Y.Z>")
	case "cut":
		fmt.Fprintln(out, "usage: cleo release cut --version <vX.Y.Z>")
	case "publish":
		fmt.Fprintln(out, "usage: cleo release publish --version <vX.Y.Z> [--draft|--final] [--no-notes]")
	case "verify":
		fmt.Fprintln(out, "usage: cleo release verify --version <vX.Y.Z>")
	case "go":
		fmt.Fprintln(out, "usage: cleo release go <plan|cut|publish|verify> --version <vX.Y.Z>")
	default:
		return false
	}
	return true
}

func PrintReleaseGo(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo release go <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  plan --version <vX.Y.Z>")
	fmt.Fprintln(out, "  cut --version <vX.Y.Z>")
	fmt.Fprintln(out, "  publish --version <vX.Y.Z> [--draft|--final] [--no-notes]")
	fmt.Fprintln(out, "  verify --version <vX.Y.Z>")
}
