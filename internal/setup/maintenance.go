package setup

import (
	"fmt"
	"io"
	"os"

	"github.com/cafaye/cleo/internal/qacatalog"
	"gopkg.in/yaml.v3"
)

func ApplyPostUpdateMigrations(out io.Writer) error {
	if _, err := os.Stat("cleo.yml"); err != nil {
		return nil
	}
	changed, err := ensureConfigDefaults("cleo.yml")
	if err != nil {
		return err
	}
	if changed && out != nil {
		fmt.Fprintln(out, "Updated cleo.yml with missing defaults.")
	}
	if err := qacatalog.EnsureQAKit("."); err != nil {
		return err
	}
	if out != nil {
		fmt.Fprintln(out, "Ensured QA kit assets.")
	}
	return nil
}

func ensureConfigDefaults(path string) (bool, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	cfg := map[string]any{}
	if err := yaml.Unmarshal(body, &cfg); err != nil {
		return false, err
	}
	changed := false
	qa := ensureMap(cfg, "qa", &changed)
	if _, ok := qa["actors_dir"]; !ok {
		qa["actors_dir"] = ".cleo/qa/actors"
		changed = true
	}
	if _, ok := qa["evidence_dir"]; !ok {
		qa["evidence_dir"] = ".cleo/evidence"
		changed = true
	}
	manual := ensureMap(qa, "manual", &changed)
	if _, ok := manual["enabled"]; !ok {
		manual["enabled"] = true
		changed = true
	}
	if !changed {
		return false, nil
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return false, err
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func ensureMap(parent map[string]any, key string, changed *bool) map[string]any {
	if raw, ok := parent[key]; ok {
		if m, ok := raw.(map[string]any); ok {
			return m
		}
	}
	m := map[string]any{}
	parent[key] = m
	*changed = true
	return m
}
