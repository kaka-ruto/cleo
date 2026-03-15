package setup

import (
	"fmt"
	"io"

	"github.com/cafaye/cleo/internal/qacatalog"
)

func ApplyPostUpdateMigrations(out io.Writer) error {
	if err := qacatalog.EnsureQAKit("."); err != nil {
		return err
	}
	if out != nil {
		fmt.Fprintln(out, "Ensured QA kit assets.")
	}
	return nil
}
