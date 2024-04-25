package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

func (r *Renderer) renderParam(node *ast.Param, w *render.BufferedWriter) error {
	w.Print(node.ParamName)
	w.Print(" ")

	return r.renderNode(node.ParamTypeDecl, w)
}
