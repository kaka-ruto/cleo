package pr

import (
	"strings"
	"testing"
)

func TestParseCommandsParsesList(t *testing.T) {
	body := `x
<!-- post-merge-commands:start -->
- ` + "`bin/kamal runner-primary -- \"puts 1\"`" + `
- ` + "`None`" + `
<!-- post-merge-commands:end -->`
	cmds, err := parseCommands(body, "<!-- post-merge-commands:start -->", "<!-- post-merge-commands:end -->", true)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(cmds) != 1 || !strings.HasPrefix(cmds[0], "bin/kamal") {
		t.Fatalf("unexpected commands: %#v", cmds)
	}
}

func TestParseCommandsFailsWithoutMarkers(t *testing.T) {
	_, err := parseCommands("missing", "a", "b", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAllowedAndDenied(t *testing.T) {
	if !allowed("bin/kamal logs", []string{"bin/kamal"}) {
		t.Fatal("expected allowed")
	}
	if denied("bin/rails test", []string{"bin/rails"}) == false {
		t.Fatal("expected denied")
	}
}
