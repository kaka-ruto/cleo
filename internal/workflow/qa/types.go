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
	Start(source string, ref string, goals string) (int64, error)
	LogIssue(sessionID int64, title string, details string, severity string) (int64, bool, error)
	Finish(sessionID int64, verdict string) error
	Report(sessionID int64) (string, error)
	Plan(sessionID int64, acFile string) (string, error)
	Run(sessionID int64, acFile string) (string, error)
	Doctor(acFile string) (string, error)
}
