package setup

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kaka-ruto/cleo/internal/qacatalog"
	"github.com/kaka-ruto/cleo/internal/skills"
)

func ApplyPostUpdateMigrations(out io.Writer) error {
	if err := qacatalog.EnsureQAKit("."); err != nil {
		return err
	}
	if out != nil {
		fmt.Fprintln(out, "Ensured QA kit assets.")
	}
	created, err := ensureBuiltinSkill(".", "cleo")
	if err != nil {
		return err
	}
	if out != nil {
		if created {
			fmt.Fprintln(out, "Installed builtin cleo skill at .agents/skills/cleo/SKILL.md.")
		} else {
			fmt.Fprintln(out, "Found existing .agents/skills/cleo/SKILL.md; leaving it unchanged.")
		}
	}
	return nil
}

func ensureBuiltinSkill(repoRoot string, name string) (bool, error) {
	target := filepath.Join(repoRoot, ".agents", "skills", name, "SKILL.md")
	if _, err := os.Stat(target); err == nil {
		return false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	body, err := skills.ReadBuiltin(name)
	if err != nil {
		return false, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return false, err
	}
	if err := os.WriteFile(target, body, 0o644); err != nil {
		return false, err
	}
	return true, nil
}
