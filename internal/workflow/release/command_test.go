package release

import "testing"

type fakeActions struct {
	publishedVersion string
	draft            bool
	notes            bool
}

func (f *fakeActions) CheckGitClean() error              { return nil }
func (f *fakeActions) EnsureReleaseMissing(string) error { return nil }
func (f *fakeActions) Cut(string) error                  { return nil }
func (f *fakeActions) Verify(string) error               { return nil }
func (f *fakeActions) Publish(v string, d bool, n bool) error {
	f.publishedVersion = v
	f.draft = d
	f.notes = n
	return nil
}

func TestExecutePublish(t *testing.T) {
	f := &fakeActions{}
	cmd := New(f, Options{TagPrefix: "v", DefaultDraft: false, GenerateNotes: true})
	if err := cmd.Execute("publish", []string{"--version", "v1.0.0", "--draft"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.publishedVersion != "v1.0.0" || !f.draft || !f.notes {
		t.Fatalf("unexpected publish call: version=%s draft=%v notes=%v", f.publishedVersion, f.draft, f.notes)
	}
}

func TestPrintOutcome(t *testing.T) {
	printOutcome(Plan{Name: "plan", Version: "v1.0.0"})
	printOutcome(Plan{Name: "cut", Version: "v1.0.0"})
	printOutcome(Plan{Name: "publish", Version: "v1.0.0"})
	printOutcome(Plan{Name: "verify", Version: "v1.0.0"})
}
