package pr

import "testing"

type fakeActions struct {
	statusPR string
}

func (f *fakeActions) Status(pr string) error { f.statusPR = pr; return nil }
func (f *fakeActions) Gate(string) error      { return nil }
func (f *fakeActions) Checks(string) error    { return nil }
func (f *fakeActions) Watch(string) error     { return nil }
func (f *fakeActions) Doctor() error          { return nil }
func (f *fakeActions) Run(string, bool) error { return nil }
func (f *fakeActions) Merge(string, bool, bool, bool, bool) error {
	return nil
}
func (f *fakeActions) Batch(int, bool, bool, bool) error { return nil }
func (f *fakeActions) Rebase(string) error               { return nil }
func (f *fakeActions) Retarget(string, string) error     { return nil }
func (f *fakeActions) Create(string, string, string, string, string, string, string, string, string, []string, bool) error {
	return nil
}

func TestExecuteStatus(t *testing.T) {
	fake := &fakeActions{}
	cmd := New(fake)
	if err := cmd.Execute("status", []string{"77"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.statusPR != "77" {
		t.Fatalf("expected status call for PR 77, got %q", fake.statusPR)
	}
}
