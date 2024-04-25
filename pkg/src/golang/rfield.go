package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/golang/validate"
	"github.com/worldiety/macro/pkg/src/render"
	"strconv"
)

func (r *Renderer) renderField(node *ast.Field, w *render.BufferedWriter) error {
	r.writeCommentNode(w, false, node.Identifier(), node.Comment())

	if err := validate.ExportedIdentifier(node.Visibility(), node.Identifier()); err != nil {
		return err
	}

	w.Printf(node.FieldName)
	w.Printf(" ")
	if err := r.renderTypeDecl(node.TypeDecl(), w); err != nil {
		return err
	}

	// try to translate annotations on fields into go struct tags
	if len(node.Annotations()) > 0 {
		w.Printf(" `")
		for i, annotationNode := range node.Annotations() {
			w.Printf(string(annotationNode.Identifier()))
			w.Printf(":")
			// by definition the go renderer ever uses the empty value as is
			v := fmt.Sprintf("%v", annotationNode.GetLiteral(""))
			w.Printf(strconv.Quote(v))

			if i < len(node.Annotations())-1 {
				w.Printf(" ") // separator between tag fields is the white space
			}
		}
		w.Printf("`")
	}

	w.Printf("\n")

	return nil
}
