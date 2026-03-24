package registry

import "testing"

func TestUpsertResolveRemoveCustom(t *testing.T) {
	home := t.TempDir()
	d := Definition{
		Name:        "team",
		Description: "Team skills",
		Repo:        "acme/skills",
		Ref:         "main",
		Path:        "skills",
	}
	if err := UpsertCustom(home, d); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	got, ok, err := ResolveDefinition(home, "team")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !ok {
		t.Fatalf("expected registry to exist")
	}
	if got.Repo != "acme/skills" || got.Path != "skills" {
		t.Fatalf("unexpected registry: %#v", got)
	}
	removed, err := RemoveCustom(home, "team")
	if err != nil {
		t.Fatalf("remove: %v", err)
	}
	if !removed {
		t.Fatalf("expected remove true")
	}
	_, ok, err = ResolveDefinition(home, "team")
	if err != nil {
		t.Fatalf("resolve after remove: %v", err)
	}
	if ok {
		t.Fatalf("expected registry removed")
	}
}
