package qaassets

type Profile struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Runbooks    []string `yaml:"runbooks"`
}

type Runbook struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Checks      []RunCheck `yaml:"checks"`
}

type RunCheck struct {
	ID             string `yaml:"id"`
	Title          string `yaml:"title"`
	Command        string `yaml:"command"`
	Severity       string `yaml:"severity"`
	FailureTitle   string `yaml:"failure_title"`
	FailureDetails string `yaml:"failure_details"`
}

type Environment struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Vars        map[string]string `yaml:"vars"`
}
