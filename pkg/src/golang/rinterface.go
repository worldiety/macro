package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/golang/validate"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderStruct emits a struct type.
func (r *Renderer) renderInterface(node *ast.Interface, w *render.BufferedWriter) error {
	r.writeCommentNode(w, false, node.Identifier(), node.Comment())

	if node.TypeName == "" {
		w.Printf(" %s interface {\n", node.Identifier())
	} else {
		if err := validate.ExportedIdentifier(node.Visibility(), node.Identifier()); err != nil {
			return err
		}

		w.Printf(" type %s interface {\n", node.Identifier())
	}

	/*
		for _, typeNode := range node.Types() {
			if err := r.renderType(typeNode, w); err != nil {
				return err
			}
		}*/

	for _, fun := range node.Methods() {
		funComment := r.renderFuncComment(fun)
		if err := r.renderFunc(fun, w); err != nil {
			return fmt.Errorf("cannot render func '%s': %w", fun.Identifier(), err)
		}

		// I like a new line after a func but be more compact without comment
		if funComment != "" {
			w.Printf("\n")
		}
	}

	for _, decl := range node.Embedded {
		if err := r.renderTypeDecl(decl, w); err != nil {
			return fmt.Errorf("cannot render embedded decl '%s': %w", decl.String(), err)
		}
	}

	w.Printf("\n}\n")

	return nil
}
