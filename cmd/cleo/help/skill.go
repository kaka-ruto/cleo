package help

import (
	"fmt"
	"io"
)

func PrintSkill(out io.Writer) {
	fmt.Fprintln(out, "usage: cleo skill <command>")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  list                             List available skills and their source")
	fmt.Fprintln(out, "  use <name>                       Print resolved SKILL.md for immediate agent use")
	fmt.Fprintln(out, "  customize <name>                 Create project override at .cleo/skills/<name>/SKILL.md")
	fmt.Fprintln(out, "  check [name]                     Validate one or all skills")
	fmt.Fprintln(out, "  help [command]")
}

func PrintSkillCommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "list":
		fmt.Fprintln(out, "usage: cleo skill list")
	case "use":
		fmt.Fprintln(out, "usage: cleo skill use <name>")
	case "customize":
		fmt.Fprintln(out, "usage: cleo skill customize <name>")
	case "check":
		fmt.Fprintln(out, "usage: cleo skill check [name]")
	default:
		return false
	}
	return true
}
