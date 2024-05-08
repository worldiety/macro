package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
)

func (r *RFile) renderStructClass(def *wdl.Struct, w *render.Writer) error {

	r.parent.writeCommentNode(w, false, "", 0, wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(def.Comment().Lines())
	}))
	w.Printf("export class %s", tsUpperName(def))
	if len(def.TypeParams()) > 0 {
		w.Printf("<")
		for i, resolvedType := range def.TypeParams() {
			w.Printf(resolvedType.Name().String())
			if i != len(def.TypeParams())-1 {
				w.Printf(", ")
			}
		}
		w.Printf(">")
	}
	w.Printf(" {\n")
	for _, field := range def.Fields() {
		if field.Visibility() != wdl.Public {
			continue
		}
		w.Printf("    private _%s : %s;\n", fieldName(field), r.TsType(field.TypeDef()))
	}

	// emit getter and setter for _ properties
	for _, field := range def.Fields() {
		if field.Visibility() != wdl.Public {
			continue
		}
		w.Printf("    get %s(): %s{\n", tsLowerName(field), r.TsType(field.TypeDef()))
		w.Printf("        return this._%s;\n", fieldName(field))
		w.Printf("    }\n")

		w.Printf("    set %s(value: %s){\n", tsLowerName(field), r.TsType(field.TypeDef()))
		w.Printf("        this._%s = value;\n", fieldName(field))
		w.Printf("    }\n")
	}

	for _, f := range def.Methods() {
		if err := r.renderFunc(f, w); err != nil {
			return err
		}
		w.Print("\n")
	}

	w.Printf("}\n\n")

	return nil
}

func fieldName(f *wdl.Field) string {
	if v, ok := f.Tags()["json"]; ok {
		return v
	}

	return f.Name().String() // without tag, in go the name will get serialized as uppercase, we must not change that
}
