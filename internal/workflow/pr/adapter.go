package pr

import corepr "github.com/cafaye/cleo/internal/pr"

type Adapter struct {
	service *corepr.Service
}

func NewAdapter(service *corepr.Service) *Adapter {
	return &Adapter{service: service}
}

func (a *Adapter) Status(pr string) error { return a.service.Status(pr) }
func (a *Adapter) Gate(pr string) error   { return a.service.Gate(pr) }
func (a *Adapter) Checks(pr string) error { return a.service.Checks(pr) }
func (a *Adapter) Watch(ref string) error { return a.service.Watch(ref) }
func (a *Adapter) Doctor() error          { return a.service.Doctor() }
func (a *Adapter) Run(pr string, dry bool) error {
	return a.service.Run(pr, dry)
}
func (a *Adapter) Merge(pr string, noWatch bool, noRun bool, noRebase bool, deleteBranch bool) error {
	return a.service.Merge(pr, noWatch, noRun, noRebase, deleteBranch)
}
func (a *Adapter) Batch(start int, noWatch bool, noRun bool, noRebase bool) error {
	return a.service.Batch(start, noWatch, noRun, noRebase)
}
func (a *Adapter) Rebase(pr string) error { return a.service.Rebase(pr) }
func (a *Adapter) Retarget(pr, base string) error {
	return a.service.Retarget(pr, base)
}
func (a *Adapter) Create(title, summary, why, what, test, risk, rollback, owner, ac string, cmds []string, draft bool) error {
	return a.service.Create(title, summary, why, what, test, risk, rollback, owner, ac, cmds, draft)
}
