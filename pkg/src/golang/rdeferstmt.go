package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

func (r *Renderer) renderDeferStmt(node *ast.DeferStmt, w *render.BufferedWriter) error {
	w.Print("defer ")

	return r.renderNode(node.CallExpr, w)
}
