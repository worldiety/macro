package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderBasicLit emits a basic literal like a string or float.
func (r *Renderer) renderBasicLit(node *ast.BasicLit, w *render.BufferedWriter) error {
	w.Printf(node.Val)

	return nil
}
