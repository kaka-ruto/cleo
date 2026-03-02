package release

func Execute(a Actions, in Input, opts Options) (Result, error) {
	switch in.Name {
	case "list":
		limit, err := intFlag(in.Args, "--limit", 20)
		if err != nil {
			return Result{}, err
		}
		return Result{Name: "list"}, a.List(limit)
	case "latest":
		return Result{Name: "latest"}, a.Latest()
	}
	version, err := versionFromArgs(in.Args, opts.TagPrefix)
	if err != nil {
		return Result{}, err
	}
	switch in.Name {
	case "plan":
		if err := a.CheckGitClean(); err != nil {
			return Result{}, err
		}
		if err := a.EnsureReleaseMissing(version); err != nil {
			return Result{}, err
		}
		return Result{Name: "plan", Version: version}, nil
	case "cut":
		if err := a.CheckGitClean(); err != nil {
			return Result{}, err
		}
		if err := a.EnsureReleaseMissing(version); err != nil {
			return Result{}, err
		}
		return Result{Name: "cut", Version: version}, a.Cut(version)
	case "publish":
		if err := a.EnsureReleaseMissing(version); err != nil {
			return Result{}, err
		}
		draft := opts.DefaultDraft
		if hasFlag(in.Args, "--draft") {
			draft = true
		}
		if hasFlag(in.Args, "--final") {
			draft = false
		}
		generateNotes := opts.GenerateNotes
		if hasFlag(in.Args, "--no-notes") {
			generateNotes = false
		}
		return Result{Name: "publish", Version: version}, a.Publish(version, draft, generateNotes)
	case "verify":
		return Result{Name: "verify", Version: version}, a.Verify(version)
	default:
		return Result{}, BuildUnknownError(in.Name)
	}
}
