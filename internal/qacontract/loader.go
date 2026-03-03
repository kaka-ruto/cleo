package qacontract

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Document, error) {
	body, err := os.ReadFile(strings.TrimSpace(path))
	if err != nil {
		return Document{}, fmt.Errorf("read AC file: %w", err)
	}
	return LoadBytes(body)
}

func LoadString(raw string) (Document, error) {
	return LoadBytes([]byte(raw))
}

func LoadBytes(body []byte) (Document, error) {
	var doc Document
	dec := yaml.NewDecoder(bytes.NewReader(body))
	dec.KnownFields(true)
	if err := dec.Decode(&doc); err != nil {
		return Document{}, fmt.Errorf("parse AC: %w", err)
	}
	if doc.Version == 0 {
		doc.Version = 1
	}
	if doc.Version != 1 {
		return Document{}, fmt.Errorf("ac version must be 1")
	}
	if len(doc.Criteria) == 0 {
		return Document{}, fmt.Errorf("ac criteria are required")
	}
	for _, c := range doc.Criteria {
		c = normalizeCriterion(c)
		if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Title) == "" {
			return Document{}, fmt.Errorf("each criterion requires id and title")
		}
		if len(c.Actors) == 0 {
			return Document{}, fmt.Errorf("criterion %q requires actors", c.ID)
		}
		if strings.TrimSpace(c.Surface) == "" {
			return Document{}, fmt.Errorf("criterion %q requires surface", c.ID)
		}
		if strings.TrimSpace(c.Given) == "" || strings.TrimSpace(c.When) == "" || len(c.Then) == 0 {
			return Document{}, fmt.Errorf("criterion %q requires given, when, and then", c.ID)
		}
		for _, expected := range c.Then {
			if strings.TrimSpace(expected) == "" {
				return Document{}, fmt.Errorf("criterion %q has empty then entry", c.ID)
			}
		}
	}
	for i := range doc.Criteria {
		doc.Criteria[i] = normalizeCriterion(doc.Criteria[i])
	}
	return doc, nil
}

func normalizeCriterion(c Criterion) Criterion {
	if strings.TrimSpace(c.Surface) == "" {
		c.Surface = strings.TrimSpace(c.Execution.Surface)
	}
	if strings.TrimSpace(c.Environment) == "" {
		c.Environment = strings.TrimSpace(c.Execution.Environment)
	}
	if strings.TrimSpace(c.Given) == "" {
		c.Given = strings.TrimSpace(c.Acceptance.Goal)
	}
	if strings.TrimSpace(c.When) == "" {
		c.When = "the behavior is exercised"
	}
	if len(c.Then) == 0 && strings.TrimSpace(c.Acceptance.ExpectedResult) != "" {
		c.Then = []string{strings.TrimSpace(c.Acceptance.ExpectedResult)}
	}
	return c
}
