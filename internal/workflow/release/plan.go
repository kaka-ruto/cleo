package release

import "fmt"

func BuildPlan(in Input, opts Options) (Plan, error) {
	switch in.Name {
	case "list":
		return Plan{Name: in.Name, Description: "List releases", ReadOnly: true}, nil
	case "latest":
		return Plan{Name: in.Name, Description: "Show latest release", ReadOnly: true}, nil
	case "plan":
		v, err := versionFromArgs(in.Args, opts.TagPrefix)
		if err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Version: v, Description: "Validate release preconditions", ReadOnly: true}, nil
	case "cut":
		v, err := versionFromArgs(in.Args, opts.TagPrefix)
		if err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Version: v, Description: "Create and push release tag"}, nil
	case "publish":
		v, err := versionFromArgs(in.Args, opts.TagPrefix)
		if err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Version: v, Description: "Publish GitHub release"}, nil
	case "verify":
		v, err := versionFromArgs(in.Args, opts.TagPrefix)
		if err != nil {
			return Plan{}, err
		}
		return Plan{Name: in.Name, Version: v, Description: "Verify published release", ReadOnly: true}, nil
	default:
		return Plan{}, fmt.Errorf("unknown release command: %s", in.Name)
	}
}
