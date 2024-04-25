package ast

// A DeferStmt represents a deferred call expression.
//  Go:
//    defer resource.Close()
type DeferStmt struct {
	CallExpr Node
	Obj
}

func NewDeferStmt(expr Node) *DeferStmt {
	n := &DeferStmt{}
	n.SetCallExpr(expr)
	return n
}

func (n *DeferStmt) SetCallExpr(expr Node) *DeferStmt {
	assertNotAttached(expr)
	assertSettableParent(expr).SetParent(n)
	n.CallExpr = expr

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *DeferStmt) Children() []Node {
	return []Node{n.CallExpr}
}
