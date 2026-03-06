package cost

type Command struct{}

func New() *Command {
	return &Command{}
}

func (c *Command) Execute(name string, args []string) error {
	in := Input{Name: name, Args: args}
	plan, err := BuildPlan(in)
	if err != nil {
		return err
	}
	result, err := Execute(in)
	if err != nil {
		return err
	}
	_ = Verify(plan, result)
	return nil
}
