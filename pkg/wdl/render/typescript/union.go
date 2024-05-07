package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
)

func (r *RFile) renderUnion(def *wdl.Union, w *render.Writer) error {

	r.parent.writeCommentNode(w, false, "", 0, wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(def.Comment())
	}))
	w.Printf("export type %s = \n", tsUpperName(def))

	for _, t := range def.Types() {
		decl := r.TsType(t)
		w.Printf("| %s\n", decl)
	}

	return nil
}
