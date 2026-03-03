package qacatalog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	ActorsDir string
}

func (l Loader) LoadActor(name string) (Actor, error) {
	var out Actor
	path := filepath.Join(strings.TrimSpace(l.ActorsDir), strings.TrimSpace(name)+".yml")
	if err := decodeFile(path, &out); err != nil {
		return Actor{}, fmt.Errorf("load actor %q: %w", name, err)
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
