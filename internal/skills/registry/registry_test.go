package registry

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestListSkills(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/acme/skills/contents/skills", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ref") != "main" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
{"name":"alpha","path":"skills/alpha","type":"dir"},
{"name":"README.md","path":"skills/README.md","type":"file"},
{"name":"beta","path":"skills/beta","type":"dir"}]`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := Client{HTTP: srv.Client(), APIBase: srv.URL, TarBase: srv.URL}
	rows, err := c.ListSkills(Definition{Repo: "acme/skills", Ref: "main", Path: "skills"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rows) != 2 || rows[0].Name != "alpha" || rows[1].Name != "beta" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}

func TestInstallSkill(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/acme/skills/tar.gz/main", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		_, _ = w.Write(makeTarball(t, map[string]string{
			"skills-main/skills/alpha/SKILL.md":     "---\nname: alpha\ndescription: x\n---\n# Alpha\n",
			"skills-main/skills/alpha/scripts/a.sh": "echo hi\n",
			"skills-main/skills/beta/SKILL.md":      "---\nname: beta\ndescription: y\n---\n# Beta\n",
		}))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	root := t.TempDir()
	c := Client{HTTP: srv.Client(), APIBase: srv.URL, TarBase: srv.URL}
	path, err := c.InstallSkill(Definition{Repo: "acme/skills", Ref: "main", Path: "skills"}, "alpha", root, false)
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected skill path: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "alpha", "scripts", "a.sh")); err != nil {
		t.Fatalf("expected bundled script: %v", err)
	}
}

func makeTarball(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	tw := tar.NewWriter(gz)
	for name, body := range files {
		hdr := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(body))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("write hdr: %v", err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatalf("write body: %v", err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return b.Bytes()
}
