package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Store struct {
	Version    int          `yaml:"version"`
	Registries []Definition `yaml:"registries"`
}

func ConfigPath(home string) string {
	return filepath.Join(home, ".agents", "skills", "registries.yml")
}

func LoadCustom(home string) ([]Definition, error) {
	path := ConfigPath(home)
	body, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var s Store
	if err := yaml.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("parse registry config: %w", err)
	}
	out := make([]Definition, 0, len(s.Registries))
	for _, d := range s.Registries {
		if err := validateDefinition(d); err != nil {
			return nil, fmt.Errorf("invalid registry %q: %w", d.Name, err)
		}
		out = append(out, normalizeDefinition(d))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func SaveCustom(home string, defs []Definition) error {
	path := ConfigPath(home)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	sort.Slice(defs, func(i, j int) bool { return defs[i].Name < defs[j].Name })
	s := Store{Version: 1, Registries: defs}
	body, err := yaml.Marshal(&s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func UpsertCustom(home string, d Definition) error {
	if err := validateDefinition(d); err != nil {
		return err
	}
	d = normalizeDefinition(d)
	defs, err := LoadCustom(home)
	if err != nil {
		return err
	}
	for _, b := range Builtins() {
		if b.Name == d.Name {
			return fmt.Errorf("cannot override builtin registry: %s", d.Name)
		}
	}
	replaced := false
	for i := range defs {
		if defs[i].Name == d.Name {
			defs[i] = d
			replaced = true
			break
		}
	}
	if !replaced {
		defs = append(defs, d)
	}
	return SaveCustom(home, defs)
}

func RemoveCustom(home string, name string) (bool, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return false, fmt.Errorf("registry name is required")
	}
	defs, err := LoadCustom(home)
	if err != nil {
		return false, err
	}
	out := make([]Definition, 0, len(defs))
	removed := false
	for _, d := range defs {
		if d.Name == name {
			removed = true
			continue
		}
		out = append(out, d)
	}
	if !removed {
		return false, nil
	}
	return true, SaveCustom(home, out)
}

func ResolveDefinition(home string, name string) (Definition, bool, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, d := range Builtins() {
		if d.Name == name {
			return d, true, nil
		}
	}
	custom, err := LoadCustom(home)
	if err != nil {
		return Definition{}, false, err
	}
	for _, d := range custom {
		if d.Name == name {
			return d, true, nil
		}
	}
	return Definition{}, false, nil
}

func AllDefinitions(home string) ([]Definition, error) {
	all := append([]Definition{}, Builtins()...)
	custom, err := LoadCustom(home)
	if err != nil {
		return nil, err
	}
	all = append(all, custom...)
	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
	return all, nil
}

func validateDefinition(d Definition) error {
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(d.Repo) == "" {
		return fmt.Errorf("repo is required")
	}
	if strings.TrimSpace(d.Path) == "" {
		return fmt.Errorf("path is required")
	}
	return nil
}

func normalizeDefinition(d Definition) Definition {
	d.Name = strings.ToLower(strings.TrimSpace(d.Name))
	d.Description = strings.TrimSpace(d.Description)
	d.Repo = strings.TrimSpace(strings.TrimPrefix(d.Repo, "https://github.com/"))
	d.Path = strings.Trim(strings.TrimSpace(d.Path), "/")
	d.Ref = strings.TrimSpace(d.Ref)
	if d.Ref == "" {
		d.Ref = "main"
	}
	return d
}
