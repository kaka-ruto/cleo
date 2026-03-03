package pr

import "strings"

const (
	acStartMarker = "<!-- cleo-ac:start -->"
	acEndMarker   = "<!-- cleo-ac:end -->"
)

func renderACBlock(ac string) string {
	trimmed := strings.TrimSpace(ac)
	if trimmed == "" {
		trimmed = defaultACScaffold()
	}
	return acStartMarker + "\n" + trimmed + "\n" + acEndMarker
}

func defaultACScaffold() string {
	return `version: 1
name: Acceptance Criteria
criteria:
  - id: c1
    title: Replace with criterion title
    severity: medium
    actors: [core]
    surface: web
    environment: local
    given: Replace with setup state and actor context
    when: Replace with user/system action under test
    then:
      - Replace with observable expected outcome
    evidence_required:
      - replace_with_evidence_artifact`
}
