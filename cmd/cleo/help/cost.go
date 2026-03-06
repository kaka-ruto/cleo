package help

import (
	"fmt"
	"io"
)

func PrintCost(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo cost <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  estimate    Estimate engineering and team-loaded cost from codebase metrics")
	fmt.Fprintln(out, "  help [command]")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "examples:")
	fmt.Fprintln(out, "  cleo cost estimate")
	fmt.Fprintln(out, "  cleo cost estimate --path . --rates-source cached")
	fmt.Fprintln(out, "  cleo cost estimate --rates-source live --country Kenya")
	fmt.Fprintln(out, "  cleo cost estimate --rates-source manual --hourly-rate 160")
}

func PrintCostCommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "estimate":
		fmt.Fprintln(out, "usage: cleo cost estimate [--path <dir>] [--rates-source <cached|manual|live>] [--country <name|ISO2>] [--hourly-rate <number>]")
	default:
		return false
	}
	return true
}
