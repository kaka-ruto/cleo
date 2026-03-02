package qaassets

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	ProfilesDir     string
	RunbooksDir     string
	EnvironmentsDir string
}

func (l Loader) LoadEnvironment(name string) (Environment, error) {
	var out Environment
	path := filepath.Join(strings.TrimSpace(l.EnvironmentsDir), strings.TrimSpace(name)+".yml")
	if err := decodeFile(path, &out); err != nil {
		return Environment{}, fmt.Errorf("load environment %q: %w", name, err)
	}
	if strings.TrimSpace(out.Name) == "" {
		out.Name = strings.TrimSpace(name)
	}
	return out, nil
}

func (l Loader) LoadProfile(name string) (Profile, error) {
	var out Profile
	path := filepath.Join(strings.TrimSpace(l.ProfilesDir), strings.TrimSpace(name)+".yml")
	if err := decodeFile(path, &out); err != nil {
		return Profile{}, fmt.Errorf("load profile %q: %w", name, err)
	}
	if strings.TrimSpace(out.Name) == "" {
		out.Name = strings.TrimSpace(name)
	}
	if len(out.Runbooks) == 0 {
		return Profile{}, fmt.Errorf("profile %q must declare runbooks", name)
	}
	return out, nil
}

func (l Loader) LoadRunbook(name string) (Runbook, error) {
	var out Runbook
	path := filepath.Join(strings.TrimSpace(l.RunbooksDir), strings.TrimSpace(name)+".yml")
	if err := decodeFile(path, &out); err != nil {
		return Runbook{}, fmt.Errorf("load runbook %q: %w", name, err)
	}
	if strings.TrimSpace(out.Name) == "" {
		out.Name = strings.TrimSpace(name)
	}
	return out, nil
}

func decodeFile(path string, out any) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(body))
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}
