package typescript

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"strings"
)

func (r *RFile) renderStructIface(def *wdl.Struct, w *render.Writer) error {

	r.parent.writeCommentNode(w, false, "", 0, wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(def.Comment().Lines())
	}))
	w.Printf("export interface %s", tsUpperName(def))
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

		if field.TypeDef() == nil {
			panic(fmt.Errorf("field type %s.%s is not resolved", def.Name(), field.Name()))
		}

		fname := fieldName(field)
		if field.Comment() != nil {
			w.Printf("\n")
			r.parent.writeCommentNode(w, false, fname, 4, field.Comment())
		}

		hasOmitTag := false
		if jsonTag, ok := field.Tags()["json"]; ok {
			hasOmitTag = strings.Contains(jsonTag, "omitempty")
		}

		w.Printf("    // %s\n", field.Name())
		if constValue, ok := field.Tags()["const"]; ok {
			w.Printf("    %s: '%s'/*%s*/;\n", fieldName(field), constValue, r.TsType(field.TypeDef()))
		} else {
			jsOpt := ""
			if hasOmitTag || looksLikeOptionHack(field.TypeDef()) {
				jsOpt = "?"
			}
			if fname == tsLowerNameStr(field.Name().String()) {
				w.Printf("    %s%s: %s;\n", fieldName(field), jsOpt, r.TsType(field.TypeDef()))
			} else {
				w.Printf("    %s%s /*%s*/: %s;\n", fieldName(field), jsOpt, field.Name(), r.TsType(field.TypeDef()))
			}
		}

	}

	w.Printf("}\n\n")

	return nil
}
