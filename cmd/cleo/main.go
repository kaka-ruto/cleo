package main

import (
	"fmt"
	"os"

	"github.com/cafaye/cleo/cmd/cleo/help"
	"github.com/cafaye/cleo/internal/config"
	corepr "github.com/cafaye/cleo/internal/pr"
	"github.com/cafaye/cleo/internal/state"
	"github.com/cafaye/cleo/internal/taskstore"
	workflowpr "github.com/cafaye/cleo/internal/workflow/pr"
	workflowqa "github.com/cafaye/cleo/internal/workflow/qa"
	workflowrelease "github.com/cafaye/cleo/internal/workflow/release"
	workflowsetup "github.com/cafaye/cleo/internal/workflow/setup"
	workflowtask "github.com/cafaye/cleo/internal/workflow/task"
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
		if len(args) > 2 && args[2] == "qa" {
			help.PrintQA(os.Stdout)
			return 0
		}
		if len(args) > 2 && args[2] == "task" {
			help.PrintTask(os.Stdout)
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
		if err := workflowupdate.New(version).Execute(hasFlag(args[2:], "--non-interactive")); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	if args[1] != "pr" && args[1] != "release" && args[1] != "qa" && args[1] != "task" {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[1])
		help.PrintRoot(os.Stderr)
		return 2
	}
	if args[1] == "release" {
		return runRelease(args)
	}
	if args[1] == "qa" {
		return runQA(args)
	}
	if args[1] == "task" {
		return runTask(args)
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
		workflowrelease.NewAdapter(cfg.GitHub.Owner, cfg.GitHub.Repo, workflowrelease.NewOptions(cfg)),
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
		workflowrelease.NewAdapter(cfg.GitHub.Owner, cfg.GitHub.Repo, workflowrelease.NewOptions(cfg)),
		workflowrelease.NewOptions(cfg),
	)
	if err := cmd.Execute(args[3], args[4:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func runQA(args []string) int {
	if len(args) < 3 {
		help.PrintQA(os.Stdout)
		return 0
	}
	if args[2] == "help" || args[2] == "--help" || args[2] == "-h" {
		if len(args) > 3 && !help.PrintQACommand(os.Stdout, args[3]) {
			fmt.Fprintf(os.Stderr, "unknown qa command: %s\n\n", args[3])
			help.PrintQA(os.Stderr)
			return 2
		}
		if len(args) <= 3 {
			help.PrintQA(os.Stdout)
		}
		return 0
	}
	cfg, err := config.Load("cleo.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		return 1
	}
	path, err := state.DBPath(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "state error: %v\n", err)
		return 1
	}
	store, err := taskstore.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "state error: %v\n", err)
		return 1
	}
	defer func() { _ = store.Close() }()
	cmd := workflowqa.New(workflowqa.NewAdapter(store, state.RepoKey(cfg)))
	if err := cmd.Execute(args[2], args[3:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func runTask(args []string) int {
	if len(args) < 3 {
		help.PrintTask(os.Stdout)
		return 0
	}
	if args[2] == "help" || args[2] == "--help" || args[2] == "-h" {
		if len(args) > 3 && !help.PrintTaskCommand(os.Stdout, args[3]) {
			fmt.Fprintf(os.Stderr, "unknown task command: %s\n\n", args[3])
			help.PrintTask(os.Stderr)
			return 2
		}
		if len(args) <= 3 {
			help.PrintTask(os.Stdout)
		}
		return 0
	}
	cfg, err := config.Load("cleo.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		return 1
	}
	path, err := state.DBPath(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "state error: %v\n", err)
		return 1
	}
	store, err := taskstore.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "state error: %v\n", err)
		return 1
	}
	defer func() { _ = store.Close() }()
	cmd := workflowtask.New(workflowtask.NewAdapter(store, cfg))
	if err := cmd.Execute(args[2], args[3:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
