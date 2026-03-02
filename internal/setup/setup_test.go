package setup

import (
	"strings"
	"testing"
)

func TestDefaultConfigIncludesOwnerAndRepo(t *testing.T) {
	cfg := defaultConfig("cafaye", "cleo")
	if !strings.Contains(cfg, "owner: cafaye") {
		t.Fatal("owner missing")
	}
	if !strings.Contains(cfg, "repo: cleo") {
		t.Fatal("repo missing")
	}
	if !strings.Contains(cfg, "block_if_requested_changes") {
		t.Fatal("expected PR policy key")
	}
}
