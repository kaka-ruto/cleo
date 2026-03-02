package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cafaye/cleo/cmd/cleo/help"
	"github.com/cafaye/cleo/internal/config"
	"github.com/cafaye/cleo/internal/pr"
	"github.com/cafaye/cleo/internal/setup"
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
		help.PrintRoot(os.Stdout)
		return 0
	}
	if args[1] == "setup" {
		options := setup.Options{
			NonInteractive: hasFlag(args[2:], "--non-interactive"),
		}
		if err := setup.NewWizard(options).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	if args[1] != "pr" {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[1])
		help.PrintRoot(os.Stderr)
		return 2
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
	svc := pr.NewService(cfg)
	if err := runPR(svc, args[2], args[3:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func runPR(svc *pr.Service, cmd string, args []string) error {
	switch cmd {
	case "status":
		return runStatus(svc, args)
	case "gate":
		return runGate(svc, args)
	case "checks":
		return runChecks(svc, args)
	case "watch":
		return runWatch(svc, args)
	case "doctor":
		return runDoctor(svc, args)
	case "run":
		return runRun(svc, args)
	case "merge":
		return runMerge(svc, args)
	case "batch":
		return runBatch(svc, args)
	case "rebase":
		return runRebase(svc, args)
	case "retarget":
		return runRetarget(svc, args)
	case "create":
		return runCreate(svc, args)
	default:
		return fmt.Errorf("unknown pr command: %s", cmd)
	}
}

func runStatus(svc *pr.Service, args []string) error {
	if err := requireLen(args, 1, "cleo pr status <pr>"); err != nil {
		return err
	}
	return svc.Status(args[0])
}

func runGate(svc *pr.Service, args []string) error {
	if err := requireLen(args, 1, "cleo pr gate <pr>"); err != nil {
		return err
	}
	return svc.Gate(args[0])
}

func runChecks(svc *pr.Service, args []string) error {
	if err := requireLen(args, 1, "cleo pr checks <pr>"); err != nil {
		return err
	}
	return svc.Checks(args[0])
}

func runWatch(svc *pr.Service, args []string) error {
	if err := requireLen(args, 1, "cleo pr watch <pr|sha>"); err != nil {
		return err
	}
	return svc.Watch(args[0])
}

func runDoctor(svc *pr.Service, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("usage: cleo pr doctor")
	}
	return svc.Doctor()
}

func runRun(svc *pr.Service, args []string) error {
	if err := requireLen(args, 1, "cleo pr run <pr> [--dry]"); err != nil {
		return err
	}
	return svc.Run(args[0], hasFlag(args[1:], "--dry"))
}

func runMerge(svc *pr.Service, args []string) error {
	if err := requireLen(args, 1, "cleo pr merge <pr> [--no-watch] [--no-run] [--no-rebase] [--delete-branch]"); err != nil {
		return err
	}
	return svc.Merge(args[0], hasFlag(args[1:], "--no-watch"), hasFlag(args[1:], "--no-run"), hasFlag(args[1:], "--no-rebase"), hasFlag(args[1:], "--delete-branch"))
}

func runBatch(svc *pr.Service, args []string) error {
	start, err := parseFrom(args)
	if err != nil {
		return err
	}
	return svc.Batch(start, hasFlag(args, "--no-watch"), hasFlag(args, "--no-run"), hasFlag(args, "--no-rebase"))
}

func runRebase(svc *pr.Service, args []string) error {
	if err := requireLen(args, 1, "cleo pr rebase <pr>"); err != nil {
		return err
	}
	return svc.Rebase(args[0])
}

func runRetarget(svc *pr.Service, args []string) error {
	if err := requireLen(args, 3, "cleo pr retarget <pr> --base <branch>"); err != nil {
		return err
	}
	base := flagValue(args[1:], "--base")
	if strings.TrimSpace(base) == "" {
		return fmt.Errorf("--base is required")
	}
	return svc.Retarget(args[0], base)
}

func runCreate(svc *pr.Service, args []string) error {
	return svc.Create(
		flagValue(args, "--title"),
		flagValue(args, "--summary"),
		flagValue(args, "--why"),
		flagValue(args, "--what"),
		flagValue(args, "--test"),
		flagValue(args, "--risk"),
		flagValue(args, "--rollback"),
		flagValue(args, "--owner"),
		flagValues(args, "--cmd"),
		hasFlag(args, "--draft"),
	)
}

func parseFrom(args []string) (int, error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--from" && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				return 0, fmt.Errorf("invalid --from value: %s", args[i+1])
			}
			return n, nil
		}
	}
	return 0, nil
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func flagValue(args []string, key string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func flagValues(args []string, key string) []string {
	vals := []string{}
	for i := 0; i < len(args); i++ {
		if args[i] == key && i+1 < len(args) {
			vals = append(vals, args[i+1])
		}
	}
	return vals
}

func requireLen(args []string, min int, usage string) error {
	if len(args) < min {
		return fmt.Errorf("usage: %s", usage)
	}
	return nil
}
