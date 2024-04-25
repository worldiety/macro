package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderIfStmt emits an if statement.
func (r *Renderer) renderIfStmt(node *ast.IfStmt, w *render.BufferedWriter) error {
	if node.Init != nil {
		w.Print("if ")
		if err := r.renderNode(node.Init, w); err != nil {
			return fmt.Errorf("unable to render init: %w", err)
		}
		w.Print("; ")
	} else {
		w.Print("if ")
	}

	if err := r.renderNode(node.Cond, w); err != nil {
		return fmt.Errorf("unable to render cond: %w", err)
	}

	if err := r.renderNode(node.Body, w); err != nil {
		return fmt.Errorf("unable to render body: %w", err)
	}

	if node.Else != nil {
		if err := r.renderNode(node.Else, w); err != nil {
			return fmt.Errorf("unable to render else: %w", err)
		}
	}

	return nil
}
