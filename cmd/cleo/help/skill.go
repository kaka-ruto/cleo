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
	fmt.Fprintln(out, "  registry [list]                  List configured skill registries")
	fmt.Fprintln(out, "  registry skills <registry>       List installable skills from a registry")
	fmt.Fprintln(out, "  registry add <name> ...          Add or update a custom skill registry")
	fmt.Fprintln(out, "  registry remove <name>           Remove a custom skill registry")
	fmt.Fprintln(out, "  use <name>                       Print resolved SKILL.md for immediate agent use")
	fmt.Fprintln(out, "  customize <name>                 Create project override at .agents/skills/<name>/SKILL.md")
	fmt.Fprintln(out, "  install <name> [--global|--project] [--registry <name>] [--force]  Install one skill")
	fmt.Fprintln(out, "  uninstall <name> [--global|--project]  Remove one installed skill")
	fmt.Fprintln(out, "  sync [--global|--project]        Sync bundled skills to .agents/skills")
	fmt.Fprintln(out, "  check [name]                     Validate one or all skills")
	fmt.Fprintln(out, "  help [command]")
}

func PrintSkillCommand(out io.Writer, cmd string) bool {
	switch cmd {
	case "list":
		fmt.Fprintln(out, "usage: cleo skill list")
	case "registry":
		fmt.Fprintln(out, "usage: cleo skill registry [list]")
		fmt.Fprintln(out, "       cleo skill registry skills <registry> [--search <term>]")
		fmt.Fprintln(out, "       cleo skill registry add <name> --repo <owner/repo> --path <path> [--ref <ref>] [--description <text>]")
		fmt.Fprintln(out, "       cleo skill registry remove <name>")
	case "use":
		fmt.Fprintln(out, "usage: cleo skill use <name>")
	case "customize":
		fmt.Fprintln(out, "usage: cleo skill customize <name>")
	case "install":
		fmt.Fprintln(out, "usage: cleo skill install <name> [--global|--project] [--registry <name>] [--force]")
	case "uninstall":
		fmt.Fprintln(out, "usage: cleo skill uninstall <name> [--global|--project]")
	case "sync":
		fmt.Fprintln(out, "usage: cleo skill sync [--global|--project]")
	case "check":
		fmt.Fprintln(out, "usage: cleo skill check [name]")
	default:
		return false
	}
	return true
}
