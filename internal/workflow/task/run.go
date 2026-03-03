package task

import "fmt"

func Execute(a Actions, in Input) (Result, error) {
	switch in.Name {
	case "list":
		text, err := a.List(flagValue(in.Args, "--status"))
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "list"}, nil
	case "show":
		id, err := int64Flag(in.Args, "--id")
		if err != nil {
			return Result{}, err
		}
		text, err := a.Show(id)
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "show"}, nil
	case "claim":
		id, err := int64Flag(in.Args, "--id")
		if err != nil {
			return Result{}, err
		}
		if err := a.Claim(id); err != nil {
			return Result{}, err
		}
		fmt.Printf("Task %d claimed\n", id)
		return Result{Name: "claim"}, nil
	case "close":
		id, err := int64Flag(in.Args, "--id")
		if err != nil {
			return Result{}, err
		}
		if err := a.Close(id); err != nil {
			return Result{}, err
		}
		fmt.Printf("Task %d closed\n", id)
		return Result{Name: "close"}, nil
	case "work":
		id, err := int64Flag(in.Args, "--id")
		if err != nil {
			return Result{}, err
		}
		opts := WorkOptions{
			ForceNewBranch: hasFlag(in.Args, "--new-branch"),
			ForceInPlace:   hasFlag(in.Args, "--in-place"),
		}
		text, err := a.Work(id, opts)
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "work"}, nil
	default:
		return Result{}, BuildUnknownError(in.Name)
	}
}
