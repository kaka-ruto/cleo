package qacatalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderLoadsActor(t *testing.T) {
	root := t.TempDir()
	actors := filepath.Join(root, "actors")
	if err := os.MkdirAll(actors, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "name: buyer_web\ndescription: buyer persona\nsurfaces:\n  - web\nauth_profile: qa_buyer\n"
	if err := os.WriteFile(filepath.Join(actors, "buyer_web.yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	loader := Loader{ActorsDir: actors}
	actor, err := loader.LoadActor("buyer_web")
	if err != nil {
		t.Fatalf("LoadActor error: %v", err)
	}
	if actor.Name != "buyer_web" || actor.AuthProfile != "qa_buyer" {
		t.Fatalf("unexpected actor: %#v", actor)
	}
}
