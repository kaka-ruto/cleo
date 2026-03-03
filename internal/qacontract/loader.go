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
	var doc Document
	dec := yaml.NewDecoder(bytes.NewReader(body))
	dec.KnownFields(true)
	if err := dec.Decode(&doc); err != nil {
		return Document{}, fmt.Errorf("parse AC file: %w", err)
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
		if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Title) == "" {
			return Document{}, fmt.Errorf("each criterion requires id and title")
		}
		if strings.TrimSpace(c.Acceptance.Goal) == "" || strings.TrimSpace(c.Acceptance.ExpectedResult) == "" {
			return Document{}, fmt.Errorf("criterion %q requires acceptance.goal and acceptance.expected_result", c.ID)
		}
		if strings.TrimSpace(c.Execution.Surface) == "" {
			return Document{}, fmt.Errorf("criterion %q requires execution.surface", c.ID)
		}
		if len(c.Execution.Steps) == 0 {
			return Document{}, fmt.Errorf("criterion %q requires execution.steps", c.ID)
		}
	}
	return doc, nil
}
