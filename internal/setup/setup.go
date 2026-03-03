package setup

import (
	"fmt"

	"github.com/cafaye/cleo/internal/qacatalog"
)

func (w *Wizard) Run() error {
	w.title("Cleo Setup")
	if err := w.ensureDeps(); err != nil {
		return err
	}
	if err := w.ensurePlaywrightRuntime(); err != nil {
		return err
	}
	if err := w.ensureGitHubAuth(); err != nil {
		return err
	}
	if err := w.writeConfig(); err != nil {
		return err
	}
	if err := w.ensureQAKit(); err != nil {
		return err
	}
	if err := w.installCleo(); err != nil {
		return err
	}
	fmt.Fprintln(w.Stdout, "Setup complete. Next: cleo pr status <pr>")
	return nil
}

func (w *Wizard) ensureQAKit() error {
	fmt.Fprintln(w.Stdout, "Ensuring QA kit assets...")
	return qacatalog.EnsureQAKit(".")
}

func (w *Wizard) ensureDeps() error {
	for _, bin := range []string{"git", "gh", "gum", "node"} {
		if err := w.checkOrInstall(bin); err != nil {
			return err
		}
	}
	return nil
}
