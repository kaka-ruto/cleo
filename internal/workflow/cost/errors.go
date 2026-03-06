package cost

import "fmt"

func BuildUnknownError(name string) error {
	return fmt.Errorf("unknown cost command: %s", name)
}
