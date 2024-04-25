package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderCallExpr emits a function call.
func (r *Renderer) renderCallExpr(node *ast.CallExpr, w *render.BufferedWriter) error {
	if err := r.renderNode(node.Fun, w); err != nil {
		return fmt.Errorf("cannot render function expression: %w", err)
	}

	w.Printf("(")
	for i, n := range node.Args {
		if err := r.renderNode(n, w); err != nil {
			return fmt.Errorf("unable to render argument: %w", err)
		}

		if i < len(node.Args)-1 {
			w.Printf(", ")
		}
	}

	if node.Ellipsis {
		w.Printf("...")
	}

	w.Printf(")")

	return nil
}
