package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderBinaryExpr emits a binary expression
func (r *Renderer) renderUnaryExpr(node *ast.UnaryExpr, w *render.BufferedWriter) error {
	switch node.Op {
	case ast.OpAdd:
		w.Print("+")
	case ast.OpSub:
		w.Print("-")
	case ast.OpAnd:
		w.Print("&")
	case ast.OpNot:
		w.Print("!")
	case ast.OpInc:
	// post
	case ast.OpDec:
		// post
	default:
		panic("operator not supported: " + fmt.Sprint(node.Op))
	}

	if err := r.renderNode(node.X, w); err != nil {
		return fmt.Errorf("unable to render x: %w", err)
	}

	switch node.Op {
	case ast.OpInc:
		w.Print("++")
	case ast.OpDec:
		w.Print("--")
	}

	return nil
}
