package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// VarDecl emits a single var or multiple var block.
func (r *Renderer) renderVar(node *ast.VarDecl, w *render.BufferedWriter) error {
	if len(node.Decl) == 0 {
		return nil
	}

	if len(node.Decl) == 1 {
		r.renderAssignComment(node.Decl[0], w)
		w.Printf("var ")
		if err := r.renderNode(node.Decl[0], w); err != nil {
			return err
		}
		w.Printf("\n")

		return nil
	}

	w.Printf("var (\n")
	for _, assignment := range node.Decl {
		r.renderAssignComment(assignment, w)
		if err := r.renderNode(assignment, w); err != nil {
			return err
		}
		w.Printf("\n")
	}

	w.Printf(")\n")

	return nil
}
