package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderSelExpr emits a X.Sel expression.
func (r *Renderer) renderSelExpr(node *ast.SelExpr, w *render.BufferedWriter) error {
	if err := r.renderNode(node.X, w); err != nil {
		return fmt.Errorf("unable to render selector target: %w", err)
	}

	w.Printf(".")

	if err := r.renderIdent(node.Sel, w); err != nil {
		return fmt.Errorf("unable to render select ident: %w", err)
	}

	return nil
}
