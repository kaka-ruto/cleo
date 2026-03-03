package pr

type Input struct {
	Name string
	Args []string
}

type Plan struct {
	Name        string
	Description string
	ReadOnly    bool
}

type Result struct {
	Name string
}

type Verification struct {
	Checked bool
	Reason  string
}

type Actions interface {
	Status(pr string) error
	Gate(pr string) error
	Checks(pr string) error
	Watch(ref string) error
	Doctor() error
	Run(pr string, dry bool) error
	Merge(pr string, noWatch bool, noRun bool, noRebase bool, deleteBranch bool) error
	Batch(start int, noWatch bool, noRun bool, noRebase bool) error
	Rebase(pr string) error
	Retarget(pr, base string) error
	Create(title, summary, why, what, test, risk, rollback, owner, ac string, cmds []string, draft bool) error
}
