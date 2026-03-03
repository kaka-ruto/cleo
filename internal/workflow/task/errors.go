package task

import "fmt"

func BuildUnknownError(name string) error {
	return fmt.Errorf("unknown task command: %s", name)
}
