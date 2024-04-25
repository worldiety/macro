package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderIdent emits an identifiers name.
func (r *Renderer) renderIdent(node *ast.Ident, w *render.BufferedWriter) error {
	w.Printf(node.Name)

	return nil
}
