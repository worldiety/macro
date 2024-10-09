package golang

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
)

func (r *RFile) renderInterface(def *wdl.Interface, w *render.Writer) error {
	r.parent.writeCommentNode(w, false, "", wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(def.Comment().Lines())
	}))
	w.Printf("type %s interface {\n", goAccessorName(def))
	for _, f := range def.Methods() {
		if err := r.renderFunc(true, f, w); err != nil {
			return err
		}
		w.Print("\n")
	}
	w.Printf("}\n\n")

	return nil
}
