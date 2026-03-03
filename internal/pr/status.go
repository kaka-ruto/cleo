package pr

import (
	"fmt"
	"strings"
)

func (s *Service) Status(pr string) error {
	v, err := s.Get(pr)
	if err != nil {
		return err
	}
	fmt.Printf("PR #%d: %s\n", v.Number, v.Title)
	fmt.Printf("URL: %s\n", v.URL)
	fmt.Printf("State: %s (draft=%t, mergeable=%s)\n", v.State, v.IsDraft, v.Mergeable)
	fmt.Printf("Review decision: %s\n", valueOr(v.ReviewDecision, "UNKNOWN"))
	fmt.Printf("Base/Head: %s <- %s\n", v.BaseRefName, v.HeadRefName)
	fmt.Printf("Checks: %d\n", len(v.StatusCheckRollup))
	return nil
}

func (s *Service) Checks(pr string) error {
	v, err := s.Get(pr)
	if err != nil {
		return err
	}
	body, bodyErr := s.prBody(pr)
	if bodyErr != nil {
		fmt.Printf("Note: unable to inspect PR body markers: %v\n", bodyErr)
	}
	hasAC := prHasAC(body)
	if len(v.StatusCheckRollup) == 0 {
		fmt.Printf("No status checks reported for PR #%d yet.\n", v.Number)
		fmt.Printf("Try `cleo pr watch %d` and retry.\n", v.Number)
		return nil
	}
	e := s.evaluateChecks(v)
	fmt.Printf("Checks summary: pending=%d failed=%d total=%d\n", len(e.pending), len(e.failed), len(v.StatusCheckRollup))
	for _, c := range v.StatusCheckRollup {
		fmt.Printf("- %s status=%s conclusion=%s\n", checkLabel(c), valueOr(c.Status, "UNKNOWN"), valueOr(c.Conclusion, "UNKNOWN"))
		if c.URL != "" {
			fmt.Printf("  %s\n", c.URL)
		}
		fmt.Println("  trace: check-run-id unavailable via current rollup API; use URL for traceability.")
	}
	if len(e.pending) > 0 {
		fmt.Printf("Pending checks detected. Run `cleo pr watch %d`.\n", v.Number)
	}
	if hasAC {
		expected := qaPolicyWorkflow(body)
		if expected != "" {
			if workflowPresent(v.StatusCheckRollup, expected) {
				fmt.Printf("QA workflow: configured workflow %q check found.\n", expected)
			} else {
				fmt.Printf("QA workflow: AC markers detected but configured workflow %q check not found.\n", expected)
			}
		} else {
			matches := workflowMatches(v.StatusCheckRollup, "qa")
			if len(matches) > 0 {
				fmt.Printf("QA workflow: found QA-like checks (%s).\n", strings.Join(matches, ", "))
			} else {
				fmt.Println("QA workflow: AC markers detected but no QA-like workflow checks found.")
			}
		}
	}
	return nil
}

func (s *Service) prBody(pr string) (string, error) {
	out, err := s.gh.Run("pr", "view", pr, "--repo", s.repo(), "--json", "body", "--jq", ".body")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func prHasAC(body string) bool {
	return strings.Contains(body, "<!-- cleo-ac:start -->") && strings.Contains(body, "<!-- cleo-ac:end -->")
}

func qaPolicyWorkflow(body string) string {
	start := "<!-- cleo-qa-policy:start -->"
	end := "<!-- cleo-qa-policy:end -->"
	si := strings.Index(body, start)
	ei := strings.Index(body, end)
	if si < 0 || ei <= si {
		return ""
	}
	block := body[si+len(start) : ei]
	for _, line := range strings.Split(block, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "workflow:") {
			continue
		}
		return strings.TrimSpace(strings.TrimPrefix(trimmed, "workflow:"))
	}
	return ""
}

func workflowPresent(checks []Check, workflowName string) bool {
	target := strings.TrimSpace(workflowName)
	for _, c := range checks {
		if strings.EqualFold(strings.TrimSpace(c.WorkflowName), target) {
			return true
		}
	}
	return false
}

func workflowMatches(checks []Check, pattern string) []string {
	needle := strings.ToLower(strings.TrimSpace(pattern))
	seen := map[string]struct{}{}
	var out []string
	for _, c := range checks {
		name := strings.TrimSpace(c.WorkflowName)
		if name == "" {
			continue
		}
		if !strings.Contains(strings.ToLower(name), needle) {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}
