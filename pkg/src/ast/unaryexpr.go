package ast

// A UnaryExpr is something like a++. It is not a StarExpr and not an Assign.
type UnaryExpr struct {
	Op Operator
	X  Expr
	Obj
}

func NewUnaryExpr(x Expr, op Operator) *UnaryExpr {
	n := &UnaryExpr{
		X:  x,
		Op: op,
	}

	assertNotAttached(x)
	assertSettableParent(x).SetParent(n)

	return n
}

func (n *UnaryExpr) exprNode() {

}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *UnaryExpr) Children() []Node {
	return []Node{n.X}
}
