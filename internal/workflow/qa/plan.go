package qa

import "fmt"

func BuildPlan(in Input) (Plan, error) {
	switch in.Name {
	case "start":
		if flagValue(in.Args, "--source") == "" || flagValue(in.Args, "--ref") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa start --source <pr|request> --ref <id|text> --goals <text>")
		}
		return Plan{Name: "start", Description: "Start QA session"}, nil
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
			return Plan{}, fmt.Errorf("usage: cleo qa report --session <id>")
		}
		return Plan{Name: "report", Description: "Print QA session report", ReadOnly: true}, nil
	case "plan":
		if flagValue(in.Args, "--session") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa plan --session <id> [--env <name>] [--profiles <a,b>]")
		}
		return Plan{Name: "plan", Description: "Resolve QA profiles and runbooks", ReadOnly: true}, nil
	case "run":
		if flagValue(in.Args, "--session") == "" {
			return Plan{}, fmt.Errorf("usage: cleo qa run --session <id> [--env <name>] [--profiles <a,b>]")
		}
		return Plan{Name: "run", Description: "Execute QA runbooks"}, nil
	default:
		return Plan{}, fmt.Errorf("unknown qa command: %s", in.Name)
	}
}
