package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
)

func (r *RFile) renderDistinctType(def *wdl.DistinctType, w *render.Writer) error {

	r.parent.writeCommentNode(w, false, "", wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(def.Comment())
	}))
	decl := r.TsType(def.Underlying().AsResolvedType())
	w.Printf("export type %s = %s \n", tsUpperName(def), decl)

	for _, f := range def.Methods() {
		if err := r.renderFunc(f, w); err != nil {
			return err
		}
		w.Print("\n")
	}

	return nil
}
