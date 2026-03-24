package registry

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Definition struct {
	Name        string
	Description string
	Repo        string
	Ref         string
	Path        string
}

type Skill struct {
	Name string
	Path string
}

type Client struct {
	HTTP    *http.Client
	APIBase string
	TarBase string
}

func NewClient() Client {
	return Client{
		HTTP:    &http.Client{Timeout: 30 * time.Second},
		APIBase: "https://api.github.com",
		TarBase: "https://codeload.github.com",
	}
}

func Builtins() []Definition {
	out := []Definition{
		{
			Name:        "openai",
			Description: "OpenAI curated skills",
			Repo:        "openai/skills",
			Ref:         "main",
			Path:        "skills/.curated",
		},
		{
			Name:        "superpowers",
			Description: "obra superpowers skills",
			Repo:        "obra/superpowers",
			Ref:         "main",
			Path:        "skills",
		},
		{
			Name:        "superpowers-ruby",
			Description: "Ruby/Rails superpowers skills",
			Repo:        "lucianghinda/superpowers-ruby",
			Ref:         "main",
			Path:        "skills",
		},
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func FindDefinition(name string) (Definition, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, d := range Builtins() {
		if d.Name == name {
			return d, true
		}
	}
	return Definition{}, false
}

func (c Client) ListSkills(d Definition) ([]Skill, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/contents/%s?ref=%s", strings.TrimRight(c.APIBase, "/"), d.Repo, d.Path, d.Ref)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return nil, fmt.Errorf("registry list failed: %s: %s", res.Status, strings.TrimSpace(string(body)))
	}
	var rows []struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(res.Body).Decode(&rows); err != nil {
		return nil, err
	}
	out := make([]Skill, 0, len(rows))
	for _, r := range rows {
		if r.Type != "dir" {
			continue
		}
		out = append(out, Skill{Name: strings.ToLower(r.Name), Path: r.Path})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (c Client) InstallSkill(d Definition, skillName string, root string, force bool) (string, error) {
	skillName = strings.TrimSpace(strings.ToLower(skillName))
	if skillName == "" {
		return "", errors.New("skill name is required")
	}
	target := filepath.Join(root, skillName)
	if _, err := os.Stat(target); err == nil && !force {
		return "", fmt.Errorf("skill already installed at %s (use --force to overwrite)", target)
	}
	if force {
		if err := os.RemoveAll(target); err != nil {
			return "", fmt.Errorf("remove existing skill: %w", err)
		}
	}

	url := fmt.Sprintf("%s/%s/tar.gz/%s", strings.TrimRight(c.TarBase, "/"), d.Repo, d.Ref)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return "", fmt.Errorf("registry download failed: %s: %s", res.Status, strings.TrimSpace(string(body)))
	}

	gzr, err := gzip.NewReader(res.Body)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	want := filepath.ToSlash(filepath.Join(d.Path, skillName)) + "/"
	wrote := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		name := filepath.ToSlash(strings.TrimPrefix(hdr.Name, "./"))
		anchor := "/" + want
		idx := strings.Index(name, anchor)
		if idx < 0 {
			continue
		}
		rel := strings.TrimPrefix(name[idx+len(anchor):], "/")
		if rel == "" {
			continue
		}
		dest := filepath.Join(target, filepath.FromSlash(rel))
		if !strings.HasPrefix(filepath.Clean(dest), filepath.Clean(target)+string(os.PathSeparator)) {
			return "", fmt.Errorf("refusing path traversal in archive: %s", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return "", fmt.Errorf("create dir: %w", err)
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return "", fmt.Errorf("create dir: %w", err)
			}
			f, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
			if err != nil {
				return "", fmt.Errorf("open file: %w", err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				_ = f.Close()
				return "", fmt.Errorf("write file: %w", err)
			}
			if err := f.Close(); err != nil {
				return "", fmt.Errorf("close file: %w", err)
			}
			wrote = true
		}
	}
	if !wrote {
		_ = os.RemoveAll(target)
		return "", fmt.Errorf("skill not found in registry: %s", skillName)
	}
	skillFile := filepath.Join(target, "SKILL.md")
	if _, err := os.Stat(skillFile); err != nil {
		_ = os.RemoveAll(target)
		return "", fmt.Errorf("installed skill missing SKILL.md: %s", skillName)
	}
	return skillFile, nil
}
