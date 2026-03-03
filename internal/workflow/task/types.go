package task

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

type WorkOptions struct {
	ForceNewBranch bool
	ForceInPlace   bool
}

type Actions interface {
	List(status string) (string, error)
	Show(id int64) (string, error)
	Claim(id int64) error
	Close(id int64) error
	Work(id int64, opts WorkOptions) (string, error)
}
