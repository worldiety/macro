package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderReturnStmt emits a single or multi return statement.
func (r *Renderer) renderReturnStmt(node *ast.ReturnStmt, w *render.BufferedWriter) error {
	w.Print("return ")

	for i, result := range node.Results {
		if err := r.renderNode(result, w); err != nil {
			return fmt.Errorf("unable to render result: %w", err)
		}

		if i < len(node.Results)-1 {
			w.Print(",")
		}
	}

	w.Print("\n") // always emit a termination

	return nil
}
