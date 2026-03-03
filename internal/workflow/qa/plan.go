package qa

import "fmt"

func BuildPlan(in Input) (Plan, error) {
	switch in.Name {
	case "init":
		return Plan{Name: "init", Description: "Install QA reusable kit assets"}, nil
	case "start":
		if flagValue(in.Args, "--source") == "" || flagValue(in.Args, "--ref") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa start --source <branch|pr|request> --ref <name|id|text> --goals <text> [--ac <yaml>]")
		}
		return Plan{Name: "start", Description: "Start QA session"}, nil
	case "scaffold":
		return Plan{Name: "scaffold", Description: "Generate BDD AC scaffold", ReadOnly: true}, nil
	case "log":
		if flagValue(in.Args, "--session") == "" || flagValue(in.Args, "--title") == "" || flagValue(in.Args, "--details") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa log --session <id> --title <text> --details <text> [--severity <low|medium|high|critical>]")
		}
		return Plan{Name: "log", Description: "Log or dedupe QA finding"}, nil
	case "finish":
		if flagValue(in.Args, "--session") == "" || flagValue(in.Args, "--verdict") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa finish --session <id> --verdict <pass|fail|blocked>")
		}
		return Plan{Name: "finish", Description: "Finish QA session"}, nil
	case "report":
		if flagValue(in.Args, "--session") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa report --session <id> [--publish <pr>] [--ref <pr>]")
		}
		publish := flagValue(in.Args, "--publish")
		if publish != "" && publish != "pr" {
			return Plan{}, fmt.Errorf("--publish must be pr")
		}
		if publish == "pr" && flagValue(in.Args, "--ref") == "" {
			// allowed: resolve from session source/ref
		}
		return Plan{Name: "report", Description: "Print QA session report", ReadOnly: true}, nil
	case "plan":
		if flagValue(in.Args, "--session") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa plan --session <id>")
		}
		return Plan{Name: "plan", Description: "Validate BDD AC and resolve tool plan", ReadOnly: true}, nil
	case "run":
		if flagValue(in.Args, "--session") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa run --session <id> [--mode <auto|manual|pr>]")
		}
		mode := flagValue(in.Args, "--mode")
		if mode != "" && mode != "auto" && mode != "manual" && mode != "pr" {
			return Plan{}, fmt.Errorf("--mode must be auto|manual|pr")
		}
		return Plan{Name: "run", Description: "Run BDD AC guidance and record findings"}, nil
	case "doctor":
		if flagValue(in.Args, "--session") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa doctor --session <id>")
		}
		return Plan{Name: "doctor", Description: "Check required QA tools", ReadOnly: true}, nil
	default:
		return Plan{}, fmt.Errorf("unknown qa command: %s", in.Name)
	}
}
