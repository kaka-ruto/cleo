package cost

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
