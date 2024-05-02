package golang

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
)

func (r *RFile) renderStruct(def *wdl.Struct, w *render.Writer) error {

	r.parent.writeCommentNode(w, false, "", wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(def.Comment())
	}))
	w.Printf("type %s struct {\n", goAccessorName(def))
	for _, field := range def.Fields() {
		w.Printf("%s %s\n", goAccessorName(field), r.GoType(field.TypeDef()))
	}
	w.Printf("}\n\n")

	for _, f := range def.Methods() {
		if err := r.renderFunc(f, w); err != nil {
			return err
		}
		w.Print("\n")
	}

	return nil
}
