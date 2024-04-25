package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderForStmt emits a for statement.
func (r *Renderer) renderForStmt(node *ast.ForStmt, w *render.BufferedWriter) error {
	w.Print("for ")
	if node.Init != nil {
		if err := r.renderNode(node.Init, w); err != nil {
			return fmt.Errorf("unable to render init: %w", err)
		}
		w.Print("; ")
	}

	if node.Cond != nil {
		w.Print(" ")
		if err := r.renderNode(node.Cond, w); err != nil {
			return fmt.Errorf("unable to render cond: %w", err)
		}
	}

	if node.Post != nil {
		w.Print("; ")
		if err := r.renderNode(node.Post, w); err != nil {
			return fmt.Errorf("unable to render post: %w", err)
		}
	}

	if err := r.renderNode(node.Body, w); err != nil {
		return fmt.Errorf("unable to render body: %w", err)
	}

	return nil
}
