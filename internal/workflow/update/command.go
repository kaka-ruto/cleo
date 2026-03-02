package update

import "fmt"

type Command struct {
	updater *ReleaseUpdater
}

func New() *Command {
	return &Command{updater: NewReleaseUpdater()}
}

func (c *Command) Execute(_ bool) error {
	if err := c.updater.UpdateLatest(); err != nil {
		return fmt.Errorf("release update failed: %w", err)
	}
	return nil
}
