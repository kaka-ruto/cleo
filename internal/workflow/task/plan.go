package task

import "fmt"

func BuildPlan(in Input) (Plan, error) {
	switch in.Name {
	case "list":
		return Plan{Name: "list", Description: "List tasks", ReadOnly: true}, nil
	case "show":
		if flagValue(in.Args, "--id") == "" {
			return Plan{}, fmt.Errorf("usage: cleo task show --id <task-id>")
		}
		return Plan{Name: "show", Description: "Show task", ReadOnly: true}, nil
	case "claim":
		if flagValue(in.Args, "--id") == "" {
			return Plan{}, fmt.Errorf("usage: cleo task claim --id <task-id>")
		}
		return Plan{Name: "claim", Description: "Claim task"}, nil
	case "close":
		if flagValue(in.Args, "--id") == "" {
			return Plan{}, fmt.Errorf("usage: cleo task close --id <task-id>")
		}
		return Plan{Name: "close", Description: "Close task"}, nil
	case "work":
		if flagValue(in.Args, "--id") == "" {
			return Plan{}, fmt.Errorf("usage: cleo task work --id <task-id> [--new-branch|--in-place]")
		}
		if hasFlag(in.Args, "--new-branch") && hasFlag(in.Args, "--in-place") {
			return Plan{}, fmt.Errorf("--new-branch and --in-place cannot be used together")
		}
		return Plan{Name: "work", Description: "Start work lane for task"}, nil
	default:
		return Plan{}, fmt.Errorf("unknown task command: %s", in.Name)
	}
}
