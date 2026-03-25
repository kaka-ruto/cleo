package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kaka-ruto/cleo/internal/ghcli"
	releaseruntime "github.com/kaka-ruto/cleo/internal/workflow/release/runtime"
)

type Adapter struct {
	gh            *ghcli.Client
	repo          string
	root          string
	changelogFile string
	binaryName    string
	buildTarget   string
}

func NewAdapter(owner, repo string, opts Options) *Adapter {
	root, _ := os.Getwd()
	return &Adapter{
		gh:            ghcli.New(),
		repo:          owner + "/" + repo,
		root:          root,
		changelogFile: opts.ChangelogFile,
		binaryName:    opts.BinaryName,
		buildTarget:   opts.BuildTarget,
	}
}

func (a *Adapter) CheckGitClean() error {
	out, err := runLocal("git", "status", "--porcelain")
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("working tree is not clean")
	}
	return nil
}

func (a *Adapter) EnsureReleaseMissing(version string) error {
	_, err := a.gh.Run("release", "view", version, "--repo", a.repo)
	if err == nil {
		return fmt.Errorf("release %s already exists", version)
	}
	if strings.Contains(strings.ToLower(err.Error()), "release not found") {
		return nil
	}
	if strings.Contains(strings.ToLower(err.Error()), "http 404") {
		return nil
	}
	return err
}

func (a *Adapter) Cut(version string) error {
	if _, err := runLocal("git", "tag", version); err != nil {
		return err
	}
	_, err := runLocal("git", "push", "origin", version)
	return err
}

func (a *Adapter) Publish(version string, draft bool, generateNotes bool, notes NoteOverrides) error {
	args := []string{"release", "create", version, "--repo", a.repo, "--verify-tag"}
	if draft {
		args = append(args, "--draft")
	}
	if generateNotes {
		notesPath, err := a.writeNotesFile(version, notes)
		if err != nil {
			return err
		}
		args = append(args, "--notes-file", notesPath)
	}
	if releaseruntime.DetectGo(a.root) {
		assets, err := releaseruntime.BuildGoReleaseArtifacts(version, a.binaryName, a.buildTarget)
		if err != nil {
			return err
		}
		args = append(args, assets...)
	} else if releaseruntime.DetectRuby(a.root) {
		assets, err := releaseruntime.BuildRubyReleaseArtifacts(a.root, version)
		if err != nil {
			return err
		}
		args = append(args, assets...)
	}
	_, err := a.gh.Run(args...)
	return err
}

func (a *Adapter) writeNotesFile(version string, notes NoteOverrides) (string, error) {
	generated, err := a.generateNotes(version)
	if err != nil {
		return "", err
	}
	sections, warnings := changelogSections(a.changelogFile, version)
	for _, w := range warnings {
		fmt.Printf("Warning: %s\n", w)
	}
	sections = applyOverrides(sections, notes)
	changelogURL := fmt.Sprintf("https://github.com/%s/blob/%s/%s", a.repo, version, a.changelogFile)
	fullChangelogURL := fmt.Sprintf("https://github.com/%s/commits/%s", a.repo, version)
	body := buildReleaseNotesWithChangelog(version, generated, sections, changelogURL, fullChangelogURL)
	if err := validateReleaseNotes(body); err != nil {
		return "", err
	}
	path := filepath.Join(os.TempDir(), "cleo-release-notes-"+version+".md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func (a *Adapter) generateNotes(version string) (string, error) {
	api := fmt.Sprintf("repos/%s/releases/generate-notes", a.repo)
	out, err := a.gh.Run("api", api, "-f", "tag_name="+version)
	if err != nil {
		return "", err
	}
	var payload struct {
		Body string `json:"body"`
	}
	if err := ghcli.DecodeJSON(out, &payload); err != nil {
		return "", err
	}
	if strings.TrimSpace(payload.Body) == "" {
		return "- No PR-based notes generated.", nil
	}
	return payload.Body, nil
}

func (a *Adapter) Verify(version string) error {
	out, err := a.gh.Run("release", "view", version, "--repo", a.repo, "--json", "tagName,url,isDraft,isPrerelease,assets")
	if err != nil {
		return err
	}
	if !releaseruntime.DetectGo(a.root) {
		return nil
	}
	var payload struct {
		Assets []struct {
			Name string `json:"name"`
		} `json:"assets"`
	}
	if err := ghcli.DecodeJSON(out, &payload); err != nil {
		return err
	}
	if releaseruntime.DetectRuby(a.root) {
		names := make([]string, 0, len(payload.Assets))
		for _, asset := range payload.Assets {
			names = append(names, strings.TrimSpace(asset.Name))
		}
		return releaseruntime.VerifyRubyAssets(names)
	}
	have := map[string]bool{}
	for _, a := range payload.Assets {
		have[strings.TrimSpace(a.Name)] = true
	}
	for _, name := range releaseruntime.ExpectedGoAssetNames(version, a.binaryName) {
		if !have[name] {
			return fmt.Errorf("missing release asset: %s", name)
		}
	}
	return nil
}

func (a *Adapter) ValidateChangelog(version string) error {
	_, warnings := changelogSections(a.changelogFile, version)
	for _, w := range warnings {
		fmt.Printf("Warning: %s\n", w)
	}
	return nil
}

func applyOverrides(base ChangelogSections, o NoteOverrides) ChangelogSections {
	if strings.TrimSpace(o.Summary) != "" {
		base.Summary = o.Summary
	}
	if strings.TrimSpace(o.Highlights) != "" {
		base.Highlights = o.Highlights
	}
	if strings.TrimSpace(o.BreakingChanges) != "" {
		base.BreakingChanges = o.BreakingChanges
	}
	if strings.TrimSpace(o.MigrationNotes) != "" {
		base.MigrationNotes = o.MigrationNotes
	}
	if strings.TrimSpace(o.Verification) != "" {
		base.Verification = o.Verification
	}
	return base
}

func (a *Adapter) List(limit int) error {
	out, err := a.gh.Run("release", "list", "--repo", a.repo, "--limit", fmt.Sprintf("%d", limit))
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}

func (a *Adapter) Latest() error {
	out, err := a.gh.Run("release", "view", "--repo", a.repo, "--json", "tagName,url,isDraft,isPrerelease,publishedAt")
	if err != nil {
		return err
	}
	var payload struct {
		TagName      string `json:"tagName"`
		URL          string `json:"url"`
		IsDraft      bool   `json:"isDraft"`
		IsPrerelease bool   `json:"isPrerelease"`
		PublishedAt  string `json:"publishedAt"`
	}
	if err := ghcli.DecodeJSON(out, &payload); err != nil {
		return err
	}
	fmt.Printf("Latest: %s\n", payload.TagName)
	fmt.Printf("URL: %s\n", payload.URL)
	fmt.Printf("Draft: %t\n", payload.IsDraft)
	fmt.Printf("Prerelease: %t\n", payload.IsPrerelease)
	fmt.Printf("Published: %s\n", payload.PublishedAt)
	return nil
}

func runLocal(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
