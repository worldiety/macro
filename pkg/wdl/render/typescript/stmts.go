package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"log/slog"
)

func (r *RFile) renderStmt(stmt wdl.Statement, w *render.Writer) error {
	switch t := stmt.(type) {
	case wdl.RawStmt:
		w.Print(t)
		return nil
	default:
		slog.Error("unknown statement type in ts, ignoring in code generation: %T", t)
		return nil
	}

}
