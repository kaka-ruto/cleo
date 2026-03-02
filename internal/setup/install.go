package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

func (w *Wizard) installCleo() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	targetDir := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}
	targetPath := filepath.Join(targetDir, "cleo")
	if exePath == targetPath {
		fmt.Fprintf(w.Stdout, "cleo command is already installed at %s\n", targetPath)
		return nil
	}
	if _, err := os.Stat(targetPath); err == nil {
		overwrite, err := w.confirm(fmt.Sprintf("%s already exists. Replace it?", targetPath))
		if err != nil {
			return err
		}
		if !overwrite {
			fmt.Fprintf(w.Stdout, "Keeping existing cleo command at %s\n", targetPath)
			return nil
		}
	}
	if err := copyExecutable(exePath, targetPath); err != nil {
		return err
	}
	fmt.Fprintf(w.Stdout, "Installed cleo command to %s\n", targetPath)
	if !pathContains(targetDir) {
		fmt.Fprintf(w.Stdout, "Add %s to your PATH to use `cleo` globally.\n", targetDir)
	}
	return nil
}
