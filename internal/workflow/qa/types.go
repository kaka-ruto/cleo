package qa

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
	Init() error
	Start(source string, ref string, goals string, ac string) (int64, error)
	LogIssue(sessionID int64, title string, details string, severity string) (int64, bool, error)
	Finish(sessionID int64, verdict string) error
	Report(sessionID int64, publish string, ref string) (string, error)
	Plan(sessionID int64) (string, error)
	Run(sessionID int64, mode string) (string, error)
	Doctor(sessionID int64) (string, error)
	Scaffold(title string) (string, error)
}
