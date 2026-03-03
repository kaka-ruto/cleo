package qaaction

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cafaye/cleo/internal/qacontract"
)

type Registry struct{}

func NewRegistry() Registry {
	return Registry{}
}

func (r Registry) Validate(doc qacontract.Document) error {
	allowedSurfaces := map[string]struct{}{
		"web":    {},
		"api":    {},
		"mobile": {},
		"cli":    {},
	}
	for _, criterion := range doc.Criteria {
		surface := strings.TrimSpace(criterion.Surface)
		if _, ok := allowedSurfaces[surface]; !ok {
			return fmt.Errorf("criterion %q has unsupported surface %q", criterion.ID, criterion.Surface)
		}
		for i, evidence := range criterion.Evidence {
			if strings.TrimSpace(evidence) == "" {
				return fmt.Errorf("criterion %q evidence_required[%d] must not be empty", criterion.ID, i)
			}
		}
	}
	return nil
}

func (r Registry) ToolSummary(doc qacontract.Document) []string {
	set := map[string]struct{}{}
	for _, criterion := range doc.Criteria {
		switch strings.TrimSpace(criterion.Surface) {
		case "web":
			set["browser"] = struct{}{}
		case "api":
			set["api"] = struct{}{}
		}
		for _, evidence := range criterion.Evidence {
			lower := strings.ToLower(strings.TrimSpace(evidence))
			if strings.Contains(lower, "screenshot") || strings.Contains(lower, "video") {
				set["browser"] = struct{}{}
			}
			if strings.Contains(lower, "api") || strings.Contains(lower, "response") {
				set["api"] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(set))
	for tool := range set {
		out = append(out, tool)
	}
	sort.Strings(out)
	return out
}
