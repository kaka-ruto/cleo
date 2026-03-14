package skills

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed builtin/*/SKILL.md
var builtinFS embed.FS

type Source struct {
	Name   string
	Path   string
	Origin string
}

type Resolver struct {
	Cwd  string
	Home string
}

func NewResolver() (Resolver, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Resolver{}, fmt.Errorf("resolve cwd: %w", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Resolver{}, fmt.Errorf("resolve home: %w", err)
	}
	return Resolver{Cwd: cwd, Home: home}, nil
}

func (r Resolver) Resolve(name string) (Source, []byte, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return Source{}, nil, errors.New("skill name is required")
	}
	canonical := canonicalSkillName(name)

	for _, p := range r.overrideCandidates(name, canonical) {
		body, err := os.ReadFile(p)
		if err == nil {
			return Source{Name: canonical, Path: p, Origin: originForPath(r, p)}, body, nil
		}
	}

	body, err := builtinFS.ReadFile(filepath.ToSlash(filepath.Join("builtin", canonical, "SKILL.md")))
	if err == nil {
		return Source{Name: canonical, Path: "builtin/" + canonical + "/SKILL.md", Origin: "builtin"}, body, nil
	}
	return Source{}, nil, fmt.Errorf("skill not found: %s", name)
}

func (r Resolver) List() ([]Source, error) {
	set := map[string]Source{}
	for _, s := range builtinSources() {
		set[s.Name] = s
	}
	for _, root := range r.searchRoots() {
		rows, err := scanRoot(root)
		if err != nil {
			continue
		}
		for _, s := range rows {
			if _, ok := set[s.Name]; ok {
				continue
			}
			s.Origin = originForPath(r, s.Path)
			set[s.Name] = s
		}
	}
	out := make([]Source, 0, len(set))
	for _, s := range set {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (r Resolver) Customize(name string) (string, error) {
	src, body, err := r.Resolve(name)
	if err != nil {
		return "", err
	}
	target := filepath.Join(r.Cwd, ".agents", "skills", src.Name, "SKILL.md")
	if src.Path == target {
		return target, nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", fmt.Errorf("create skill directory: %w", err)
	}
	if err := os.WriteFile(target, body, 0o644); err != nil {
		return "", fmt.Errorf("write customized skill: %w", err)
	}
	return target, nil
}

func (r Resolver) Check(name string) ([]Source, error) {
	if strings.TrimSpace(name) != "" {
		src, body, err := r.Resolve(name)
		if err != nil {
			return nil, err
		}
		if err := validateSkill(body); err != nil {
			return nil, fmt.Errorf("%s: %w", src.Path, err)
		}
		return []Source{src}, nil
	}
	list, err := r.List()
	if err != nil {
		return nil, err
	}
	for _, s := range list {
		_, body, err := r.Resolve(s.Name)
		if err != nil {
			return nil, err
		}
		if err := validateSkill(body); err != nil {
			return nil, fmt.Errorf("%s: %w", s.Path, err)
		}
	}
	return list, nil
}

func (r Resolver) overrideCandidates(name string, canonical string) []string {
	names := uniqueNames(name, canonical)
	out := make([]string, 0, len(names)*2)
	for _, n := range names {
		out = append(out,
			filepath.Join(r.Cwd, ".agents", "skills", n, "SKILL.md"),
			filepath.Join(r.Home, ".agents", "skills", n, "SKILL.md"),
		)
	}
	return out
}

func (r Resolver) searchRoots() []string {
	return []string{
		filepath.Join(r.Cwd, ".agents", "skills"),
		filepath.Join(r.Home, ".agents", "skills"),
	}
}

func scanRoot(root string) ([]Source, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	out := []Source{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(root, e.Name(), "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			continue
		}
		out = append(out, Source{Name: strings.ToLower(e.Name()), Path: path})
	}
	return out, nil
}

func builtinSources() []Source {
	entries, err := builtinFS.ReadDir("builtin")
	if err != nil {
		return nil
	}
	out := []Source{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		out = append(out, Source{
			Name:   strings.ToLower(e.Name()),
			Path:   "builtin/" + e.Name() + "/SKILL.md",
			Origin: "builtin",
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func BuiltinList() []Source {
	return builtinSources()
}

func ReadBuiltin(name string) ([]byte, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, errors.New("skill name is required")
	}
	body, err := builtinFS.ReadFile(filepath.ToSlash(filepath.Join("builtin", name, "SKILL.md")))
	if err != nil {
		return nil, fmt.Errorf("builtin skill not found: %s", name)
	}
	return body, nil
}

func originForPath(r Resolver, path string) string {
	norm := filepath.Clean(path)
	if strings.HasPrefix(norm, filepath.Join(r.Cwd, ".agents", "skills")) {
		return "project"
	}
	if strings.HasPrefix(norm, filepath.Join(r.Home, ".agents", "skills")) {
		return "user"
	}
	if strings.HasPrefix(norm, "builtin/") {
		return "builtin"
	}
	return "unknown"
}

func ValidateForUse(body []byte, _ string) error {
	return validateSkill(body)
}

func canonicalSkillName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func uniqueNames(names ...string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(names))
	for _, n := range names {
		n = strings.ToLower(strings.TrimSpace(n))
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}

func validateSkill(body []byte) error {
	text := string(body)
	if !strings.HasPrefix(text, "---\n") {
		return errors.New("missing frontmatter start")
	}
	parts := strings.SplitN(text, "\n---\n", 2)
	if len(parts) != 2 {
		return errors.New("missing frontmatter end")
	}
	meta := struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	}{}
	if err := yaml.Unmarshal([]byte(strings.TrimPrefix(parts[0], "---\n")), &meta); err != nil {
		return fmt.Errorf("invalid frontmatter: %w", err)
	}
	if strings.TrimSpace(meta.Name) == "" {
		return errors.New("frontmatter requires name")
	}
	if strings.TrimSpace(meta.Description) == "" {
		return errors.New("frontmatter requires description")
	}
	return nil
}
