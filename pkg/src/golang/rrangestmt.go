package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderRangeStmt emits a range statement.
func (r *Renderer) renderRangeStmt(node *ast.RangeStmt, w *render.BufferedWriter) error {

	w.Print("for ")

	if node.Key != nil {
		if err := r.renderNode(node.Key, w); err != nil {
			return fmt.Errorf("unable to render key: %w", err)
		}
	}

	if node.Val != nil {
		if node.Key == nil {
			w.Print("_, ")
		} else {
			w.Print(", ")
		}

		if err := r.renderNode(node.Val, w); err != nil {
			return fmt.Errorf("unable to render val: %w", err)
		}
	}

	if node.Key != nil || node.Val != nil {
		w.Print(" := ")
	}

	w.Print(" range ")

	if err := r.renderNode(node.X, w); err != nil {
		return fmt.Errorf("unable to render range target: %w", err)
	}

	if err := r.renderNode(node.Body, w); err != nil {
		return fmt.Errorf("unable to render body: %w", err)
	}

	return nil
}
