package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderBinaryExpr emits a binary expression
func (r *Renderer) renderBinaryExpr(node *ast.BinaryExpr, w *render.BufferedWriter) error {
	if err := r.renderNode(node.X, w); err != nil {
		return fmt.Errorf("unable to render x: %w", err)
	}

	switch node.Op {
	case ast.OpAdd:
		w.Print("+")
	case ast.OpSub:
		w.Print("-")
	case ast.OpMul:
		w.Print("*")
	case ast.OpQuo:
		w.Print("/")
	case ast.OpREM:
		w.Print("%")

	case ast.OpAnd:
		w.Print("&")
	case ast.OpOr:
		w.Print("|")
	case ast.OpXOR:
		w.Print("^")
	case ast.OpShl:
		w.Print("<<")
	case ast.OpShr:
		w.Print(">>")
	case ast.OpAndNot:
		w.Print("&^")

	case ast.OpLAnd:
		w.Print("&&")
	case ast.OpLOr:
		w.Print("||")
	case ast.OpArrow:
		w.Print("<-")
	case ast.OpEqual:
		w.Print("==")
	case ast.OpLess:
		w.Print("<")
	case ast.OpGreater:
		w.Print(">")

	case ast.OpNotEqual:
		w.Print("!=")
	case ast.OpLessEqual:
		w.Print("<=")
	case ast.OpGreaterEqual:
		w.Print(">=")
	case ast.OpColon:
		w.Print(":")
	default:
		panic("operator not supported: " + fmt.Sprint(node.Op))
	}

	if err := r.renderNode(node.Y, w); err != nil {
		return fmt.Errorf("unable to render y: %w", err)
	}

	return nil
}
