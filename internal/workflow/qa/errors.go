package qa

import "fmt"

func BuildUnknownError(name string) error {
	return fmt.Errorf("unknown qa command: %s", name)
}
