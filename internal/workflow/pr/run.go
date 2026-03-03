package pr

func Execute(a Actions, in Input) (Result, error) {
	switch in.Name {
	case "status":
		return runStatus(a, in.Args)
	case "gate":
		return runGate(a, in.Args)
	case "checks":
		return runChecks(a, in.Args)
	case "watch":
		return runWatch(a, in.Args)
	case "doctor":
		return Result{Name: in.Name}, a.Doctor()
	case "run":
		return runRun(a, in.Args)
	case "merge":
		return runMerge(a, in.Args)
	case "batch":
		return runBatch(a, in.Args)
	case "rebase":
		return runRebase(a, in.Args)
	case "retarget":
		return runRetarget(a, in.Args)
	case "create":
		return runCreate(a, in.Args)
	default:
		return Result{}, BuildUnknownError(in.Name)
	}
}

func runStatus(a Actions, args []string) (Result, error) {
	return Result{Name: "status"}, a.Status(args[0])
}

func runGate(a Actions, args []string) (Result, error) {
	return Result{Name: "gate"}, a.Gate(args[0])
}

func runChecks(a Actions, args []string) (Result, error) {
	return Result{Name: "checks"}, a.Checks(args[0])
}

func runWatch(a Actions, args []string) (Result, error) {
	return Result{Name: "watch"}, a.Watch(args[0])
}

func runRun(a Actions, args []string) (Result, error) {
	return Result{Name: "run"}, a.Run(args[0], hasFlag(args[1:], "--dry"))
}

func runMerge(a Actions, args []string) (Result, error) {
	return Result{Name: "merge"}, a.Merge(args[0], hasFlag(args[1:], "--no-watch"), hasFlag(args[1:], "--no-run"), hasFlag(args[1:], "--no-rebase"), hasFlag(args[1:], "--delete-branch"))
}

func runBatch(a Actions, args []string) (Result, error) {
	start, err := parseFrom(args)
	if err != nil {
		return Result{}, err
	}
	return Result{Name: "batch"}, a.Batch(start, hasFlag(args, "--no-watch"), hasFlag(args, "--no-run"), hasFlag(args, "--no-rebase"))
}

func runRebase(a Actions, args []string) (Result, error) {
	return Result{Name: "rebase"}, a.Rebase(args[0])
}

func runRetarget(a Actions, args []string) (Result, error) {
	base, err := requireBase(args[1:])
	if err != nil {
		return Result{}, err
	}
	return Result{Name: "retarget"}, a.Retarget(args[0], base)
}

func runCreate(a Actions, args []string) (Result, error) {
	return Result{Name: "create"}, a.Create(
		flagValue(args, "--title"),
		flagValue(args, "--summary"),
		flagValue(args, "--why"),
		flagValue(args, "--what"),
		flagValue(args, "--test"),
		flagValue(args, "--risk"),
		flagValue(args, "--rollback"),
		flagValue(args, "--owner"),
		flagValue(args, "--ac"),
		flagValues(args, "--cmd"),
		hasFlag(args, "--draft"),
	)
}
