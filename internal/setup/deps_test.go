package setup

import "testing"

func TestLinuxPackageName(t *testing.T) {
	if got := linuxPackageName("go"); got != "golang-go" {
		t.Fatalf("expected golang-go, got %s", got)
	}
	if got := linuxPackageName("gh"); got != "gh" {
		t.Fatalf("expected gh, got %s", got)
	}
	if got := linuxPackageName("git"); got != "git" {
		t.Fatalf("expected git, got %s", got)
	}
}
