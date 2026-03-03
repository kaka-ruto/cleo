package task

type Command struct {
	actions Actions
}

func New(actions Actions) *Command {
	return &Command{actions: actions}
}

func (c *Command) Execute(name string, args []string) error {
	in := Input{Name: name, Args: args}
	plan, err := BuildPlan(in)
	if err != nil {
		return err
	}
	result, err := Execute(c.actions, in)
	if err != nil {
		return err
	}
	_ = Verify(plan, result)
	return nil
}
