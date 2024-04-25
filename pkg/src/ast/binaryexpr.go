package ast

import "go/token"

type Operator int

const (
	OpAdd Operator = Operator(token.ADD)
	OpSub Operator = Operator(token.SUB)
	OpMul Operator = Operator(token.MUL)
	OpQuo Operator = Operator(token.QUO)
	OpREM Operator = Operator(token.REM)

	OpAnd    Operator = Operator(token.AND)
	OpOr     Operator = Operator(token.OR)
	OpXOR    Operator = Operator(token.XOR)
	OpShl    Operator = Operator(token.SHL)
	OpShr    Operator = Operator(token.SHR)
	OpAndNot Operator = Operator(token.AND_NOT)
	OpNot    Operator = Operator(token.NOT)

	OpLAnd    Operator = Operator(token.LAND)
	OpLOr     Operator = Operator(token.LOR)
	OpArrow   Operator = Operator(token.ARROW)
	OpEqual   Operator = Operator(token.EQL)
	OpLess    Operator = Operator(token.LSS)
	OpGreater Operator = Operator(token.GTR)

	OpInc Operator = Operator(token.INC)
	OpDec Operator = Operator(token.DEC)

	OpNotEqual     Operator = Operator(token.NEQ)
	OpLessEqual    Operator = Operator(token.LEQ)
	OpGreaterEqual Operator = Operator(token.GEQ)

	OpColon Operator = Operator(token.COLON)
)

// A BinaryExpr is something like a + b. It is not a StarExpr and not an Assign.
type BinaryExpr struct {
	X  Expr
	Op Operator
	Y  Expr
	Obj
}

func NewBinaryExpr(x Expr, op Operator, y Expr) *BinaryExpr {
	n := &BinaryExpr{
		X:  x,
		Op: op,
		Y:  y,
	}

	assertNotAttached(x)
	assertSettableParent(x).SetParent(n)

	assertNotAttached(y)
	assertSettableParent(y).SetParent(n)

	return n
}

func (n *BinaryExpr) exprNode() {

}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *BinaryExpr) Children() []Node {
	return []Node{n.X, n.Y}
}
