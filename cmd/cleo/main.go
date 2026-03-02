package main

import (
	"fmt"
	"os"

	"github.com/cafaye/cleo/cmd/cleo/help"
	"github.com/cafaye/cleo/internal/config"
	corepr "github.com/cafaye/cleo/internal/pr"
	workflowpr "github.com/cafaye/cleo/internal/workflow/pr"
	workflowrelease "github.com/cafaye/cleo/internal/workflow/release"
	workflowsetup "github.com/cafaye/cleo/internal/workflow/setup"
	workflowupdate "github.com/cafaye/cleo/internal/workflow/update"
)

var version = "dev"

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) == 1 {
		help.PrintRoot(os.Stdout)
		return 0
	}
	if args[1] == "version" || args[1] == "--version" {
		fmt.Printf("cleo %s\n", version)
		return 0
	}
	if args[1] == "help" || args[1] == "--help" || args[1] == "-h" {
		if len(args) > 2 && args[2] == "pr" {
			help.PrintPR(os.Stdout)
			return 0
		}
		if len(args) > 2 && args[2] == "release" {
			help.PrintRelease(os.Stdout)
			return 0
		}
		if len(args) > 2 {
			if !help.PrintCommand(os.Stdout, args[2]) {
				fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[2])
				help.PrintRoot(os.Stderr)
				return 2
			}
			return 0
		}
		help.PrintRoot(os.Stdout)
		return 0
	}
	if args[1] == "setup" {
		if err := workflowsetup.New().Execute(hasFlag(args[2:], "--non-interactive")); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	if args[1] == "update" {
		if err := workflowupdate.New().Execute(hasFlag(args[2:], "--non-interactive")); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	if args[1] != "pr" && args[1] != "release" {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[1])
		help.PrintRoot(os.Stderr)
		return 2
	}
	if args[1] == "release" {
		return runRelease(args)
	}
	if len(args) < 3 {
		help.PrintPR(os.Stdout)
		return 0
	}
	if args[2] == "help" || args[2] == "--help" || args[2] == "-h" {
		if len(args) > 3 && !help.PrintPRCommand(os.Stdout, args[3]) {
			fmt.Fprintf(os.Stderr, "unknown pr command: %s\n\n", args[3])
			help.PrintPR(os.Stderr)
			return 2
		}
		if len(args) <= 3 {
			help.PrintPR(os.Stdout)
		}
		return 0
	}
	cfg, err := config.Load("cleo.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		return 1
	}
	cmd := workflowpr.New(workflowpr.NewAdapter(corepr.NewService(cfg)))
	if err := cmd.Execute(args[2], args[3:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func runRelease(args []string) int {
	if len(args) < 3 {
		help.PrintRelease(os.Stdout)
		return 0
	}
	if args[2] == "go" {
		return runReleaseGo(args)
	}
	if args[2] == "help" || args[2] == "--help" || args[2] == "-h" {
		if len(args) > 3 && args[3] == "go" {
			help.PrintReleaseGo(os.Stdout)
			return 0
		}
		if len(args) > 3 && !help.PrintReleaseCommand(os.Stdout, args[3]) {
			fmt.Fprintf(os.Stderr, "unknown release command: %s\n\n", args[3])
			help.PrintRelease(os.Stderr)
			return 2
		}
		if len(args) <= 3 {
			help.PrintRelease(os.Stdout)
		}
		return 0
	}
	cfg, err := config.Load("cleo.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		return 1
	}
	cmd := workflowrelease.New(
		workflowrelease.NewAdapter(cfg.GitHub.Owner, cfg.GitHub.Repo),
		workflowrelease.NewOptions(cfg),
	)
	if err := cmd.Execute(args[2], args[3:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func runReleaseGo(args []string) int {
	if len(args) < 4 {
		help.PrintReleaseGo(os.Stdout)
		return 0
	}
	if args[3] == "help" || args[3] == "--help" || args[3] == "-h" {
		help.PrintReleaseGo(os.Stdout)
		return 0
	}
	cfg, err := config.Load("cleo.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		return 1
	}
	cmd := workflowrelease.New(
		workflowrelease.NewAdapter(cfg.GitHub.Owner, cfg.GitHub.Repo),
		workflowrelease.NewOptions(cfg),
	)
	if err := cmd.Execute(args[3], args[4:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
