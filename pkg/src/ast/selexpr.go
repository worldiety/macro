package ast

// SelExpr connects an expression with an Ident. Usually found for members, fields or function invocations.
type SelExpr struct {
	X   Expr
	Sel *Ident
	Obj
}

func NewSelExpr(x Expr, sel *Ident) *SelExpr {
	n := &SelExpr{
		X:   x,
		Sel: sel,
	}

	assertNotAttached(x)
	assertSettableParent(x).SetParent(n)

	assertNotAttached(sel)
	assertSettableParent(sel).SetParent(n)

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *SelExpr) Children() []Node {
	return []Node{n.X, n.Sel}
}

func (n *SelExpr) exprNode() {

}
