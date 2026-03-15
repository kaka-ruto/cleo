package setup

import "fmt"

func (w *Wizard) writeConfig() error {
	fmt.Fprintln(w.Stdout, "No cleo.yml file is used.")
	fmt.Fprintln(w.Stdout, "Cleo infers repo context from git and applies built-in policies.")
	return nil
}
