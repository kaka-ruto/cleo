package qa

import "fmt"

func Execute(a Actions, in Input) (Result, error) {
	switch in.Name {
	case "init":
		if err := a.Init(); err != nil {
			return Result{}, err
		}
		fmt.Println("QA kit initialized.")
		return Result{Name: "init"}, nil
	case "start":
		source := flagValue(in.Args, "--source")
		ref := flagValue(in.Args, "--ref")
		goals := flagValue(in.Args, "--goals")
		id, err := a.Start(source, ref, goals, flagValue(in.Args, "--ac"))
		if err != nil {
			return Result{}, err
		}
		fmt.Printf("QA session started: %d\n", id)
		return Result{Name: "start"}, nil
	case "scaffold":
		text, err := a.Scaffold(flagValue(in.Args, "--title"))
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "scaffold"}, nil
	case "log":
		sessionID, err := int64Flag(in.Args, "--session")
		if err != nil {
			return Result{}, err
		}
		taskID, created, err := a.LogIssue(sessionID, flagValue(in.Args, "--title"), flagValue(in.Args, "--details"), flagValue(in.Args, "--severity"))
		if err != nil {
			return Result{}, err
		}
		if created {
			fmt.Printf("QA issue logged as task %d\n", taskID)
		} else {
			fmt.Printf("QA issue matched existing task %d\n", taskID)
		}
		return Result{Name: "log"}, nil
	case "finish":
		sessionID, err := int64Flag(in.Args, "--session")
		if err != nil {
			return Result{}, err
		}
		if err := a.Finish(sessionID, flagValue(in.Args, "--verdict")); err != nil {
			return Result{}, err
		}
		fmt.Printf("QA session %d finished\n", sessionID)
		return Result{Name: "finish"}, nil
	case "report":
		sessionID, err := int64Flag(in.Args, "--session")
		if err != nil {
			return Result{}, err
		}
		text, err := a.Report(sessionID, flagValue(in.Args, "--publish"), flagValue(in.Args, "--ref"))
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "report"}, nil
	case "plan":
		sessionID, err := int64Flag(in.Args, "--session")
		if err != nil {
			return Result{}, err
		}
		text, err := a.Plan(sessionID)
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "plan"}, nil
	case "run":
		sessionID, err := int64Flag(in.Args, "--session")
		if err != nil {
			return Result{}, err
		}
		mode := flagValue(in.Args, "--mode")
		if mode == "" {
			mode = "auto"
		}
		text, err := a.Run(sessionID, mode)
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "run"}, nil
	case "doctor":
		sessionID, err := int64Flag(in.Args, "--session")
		if err != nil {
			return Result{}, err
		}
		text, err := a.Doctor(sessionID)
		if err != nil {
			return Result{}, err
		}
		fmt.Println(text)
		return Result{Name: "doctor"}, nil
	default:
		return Result{}, BuildUnknownError(in.Name)
	}
}
