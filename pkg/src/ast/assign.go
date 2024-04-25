package ast

import "go/token"

type AssignKind int

const (
	AssignSimple AssignKind = AssignKind(token.ASSIGN)
	AssignDefine AssignKind = AssignKind(token.DEFINE)
	AssignAdd    AssignKind = AssignKind(token.ADD_ASSIGN)
	AssignSub    AssignKind = AssignKind(token.SUB_ASSIGN)
	AssignMul    AssignKind = AssignKind(token.MUL_ASSIGN)
	AssignRem    AssignKind = AssignKind(token.REM_ASSIGN)
)

// Assign represents things like
//  Go: a := 5 or a, b = b, a or x *= 5
//  Java: let a = 5 or b = a or x *= 5
type Assign struct {
	Lhs  []Expr
	Rhs  []Expr
	Kind AssignKind
	Obj
}

func NewSimpleAssign(leftIdent *Ident, kind AssignKind, rightLit *BasicLit) *Assign {
	return NewAssign(Exprs(leftIdent), kind, Exprs(rightLit))
}

func NewAssign(lhs []Expr, kind AssignKind, rhs []Expr) *Assign {
	n := &Assign{Kind: kind}
	for _, lh := range lhs {
		assertNotAttached(lh)
		assertSettableParent(lh).SetParent(n)
		n.Lhs = append(n.Lhs, lh)
	}

	for _, rh := range rhs {
		assertNotAttached(rh)
		assertSettableParent(rh).SetParent(n)
		n.Rhs = append(n.Rhs, rh)
	}

	return n
}

func (n *Assign) SetComment(text string) *Assign {
	n.Obj.ObjComment = NewComment(text)
	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Assign) Children() []Node {
	tmp := make([]Node, 0, len(n.Lhs)+len(n.Rhs))
	for _, arg := range n.Lhs {
		tmp = append(tmp, arg)
	}

	for _, arg := range n.Rhs {
		tmp = append(tmp, arg)
	}

	return tmp
}
