package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
)

func (r *RFile) renderStmt(stmt wdl.Statement, w *render.Writer) error {
	switch t := stmt.(type) {
	case wdl.RawStmt:
		w.Print(t)
		return nil
	default:
		return fmt.Errorf("unknown statement type: %T", t)
	}

}
