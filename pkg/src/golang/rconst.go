package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderConst emits a single or multiple const block.
func (r *Renderer) renderConst(node *ast.ConstDecl, w *render.BufferedWriter) error {
	if len(node.Assignments) == 0 {
		return nil
	}

	if len(node.Assignments) == 1 {
		r.renderAssignComment(node.Assignments[0], w)
		w.Printf("const ")
		if err := r.renderNode(node.Assignments[0], w); err != nil {
			return err
		}
		w.Printf("\n")

		return nil
	}

	w.Printf("const (\n")
	for _, assignment := range node.Assignments {
		r.renderAssignComment(assignment, w)
		if err := r.renderNode(assignment, w); err != nil {
			return err
		}
		w.Printf("\n")
	}

	w.Printf(")\n")

	return nil
}
