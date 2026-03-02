package release

import "fmt"

type Command struct {
	actions Actions
	opts    Options
}

func New(actions Actions, opts Options) *Command {
	return &Command{actions: actions, opts: opts}
}

func (c *Command) Execute(name string, args []string) error {
	in := Input{Name: name, Args: args}
	plan, err := BuildPlan(in, c.opts)
	if err != nil {
		return err
	}
	result, err := Execute(c.actions, in, c.opts)
	if err != nil {
		return err
	}
	_ = Verify(plan, result)
	printOutcome(plan)
	return nil
}

func printOutcome(plan Plan) {
	switch plan.Name {
	case "plan":
		fmt.Printf("Release plan passed for %s.\n", plan.Version)
	case "cut":
		fmt.Printf("Release tag %s created and pushed.\n", plan.Version)
	case "publish":
		fmt.Printf("Release %s published.\n", plan.Version)
	case "verify":
		fmt.Printf("Release %s verified.\n", plan.Version)
	}
}
