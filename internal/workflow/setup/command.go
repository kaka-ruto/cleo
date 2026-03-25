package setup

import basesetup "github.com/kaka-ruto/cleo/internal/setup"

type Command struct{}

func New() *Command {
	return &Command{}
}

func (c *Command) Execute(nonInteractive bool) error {
	options := basesetup.Options{NonInteractive: nonInteractive}
	return basesetup.NewWizard(options).Run()
}
