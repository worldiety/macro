package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"log/slog"
)

func (r *RFile) renderTypeDef(def wdl.TypeDef, w *render.Writer) error {

	switch d := def.(type) {
	case *wdl.Struct:
		return r.renderStruct(d, w)
	default:
		slog.Error("rendering not yet implemented", "type", fmt.Sprintf("%T", d))
	}

	return nil
}
