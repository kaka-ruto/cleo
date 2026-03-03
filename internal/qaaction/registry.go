package qaaction

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cafaye/cleo/internal/qacontract"
)

type Spec struct {
	Name           string
	Surfaces       []string
	RequiredParams []string
	Tool           string
}

type Registry struct {
	specs map[string]Spec
}

func NewRegistry() Registry {
	list := []Spec{
		{Name: "open_url", Surfaces: []string{"web"}, RequiredParams: []string{"target"}, Tool: "browser"},
		{Name: "login_form", Surfaces: []string{"web"}, RequiredParams: []string{"username_ref", "password_ref"}, Tool: "browser"},
		{Name: "add_fixture_cart", Surfaces: []string{"web", "api"}, RequiredParams: []string{"fixture"}, Tool: "api"},
		{Name: "capture_ui_value", Surfaces: []string{"web"}, RequiredParams: []string{"selector", "output_key"}, Tool: "browser"},
		{Name: "call_api", Surfaces: []string{"api", "web"}, RequiredParams: []string{"method", "url", "output_key"}, Tool: "api"},
		{Name: "assert_equal_money", Surfaces: []string{"api", "web"}, RequiredParams: []string{"left_key", "right_key"}, Tool: "assertion"},
	}
	m := make(map[string]Spec, len(list))
	for _, spec := range list {
		m[spec.Name] = spec
	}
	return Registry{specs: m}
}

func (r Registry) Validate(doc qacontract.Document) error {
	for _, criterion := range doc.Criteria {
		surface := strings.TrimSpace(criterion.Execution.Surface)
		for i, step := range criterion.Execution.Steps {
			spec, ok := r.specs[strings.TrimSpace(step.Action)]
			if !ok {
				return fmt.Errorf("criterion %q step %d uses unknown action %q", criterion.ID, i+1, step.Action)
			}
			if !supportsSurface(spec.Surfaces, surface) {
				return fmt.Errorf("criterion %q step %d action %q does not support surface %q", criterion.ID, i+1, step.Action, surface)
			}
			for _, key := range spec.RequiredParams {
				if strings.TrimSpace(step.Params[key]) == "" {
					return fmt.Errorf("criterion %q step %d action %q missing required param %q", criterion.ID, i+1, step.Action, key)
				}
			}
		}
	}
	return nil
}

func (r Registry) ToolSummary(doc qacontract.Document) []string {
	set := map[string]struct{}{}
	for _, criterion := range doc.Criteria {
		for _, step := range criterion.Execution.Steps {
			if spec, ok := r.specs[strings.TrimSpace(step.Action)]; ok {
				set[spec.Tool] = struct{}{}
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

func supportsSurface(surfaces []string, surface string) bool {
	for _, s := range surfaces {
		if s == surface {
			return true
		}
	}
	return false
}
