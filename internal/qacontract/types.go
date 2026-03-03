package qacontract

type Document struct {
	Version  int         `yaml:"version"`
	Name     string      `yaml:"name"`
	Criteria []Criterion `yaml:"criteria"`
}

type Criterion struct {
	ID         string        `yaml:"id"`
	Title      string        `yaml:"title"`
	Severity   string        `yaml:"severity"`
	Actors     []string      `yaml:"actors"`
	Acceptance Acceptance    `yaml:"acceptance"`
	Execution  ExecutionPlan `yaml:"execution"`
	Evidence   []string      `yaml:"evidence_required"`
}

type Acceptance struct {
	Goal           string `yaml:"goal"`
	ExpectedResult string `yaml:"expected_result"`
}

type ExecutionPlan struct {
	Surface       string            `yaml:"surface"`
	Environment   string            `yaml:"environment"`
	Preconditions map[string]string `yaml:"preconditions"`
	Steps         []Step            `yaml:"steps"`
}

type Step struct {
	Action string            `yaml:"action"`
	Params map[string]string `yaml:"params"`
}
