package ast

// A CallExpr declares a call to Fun with the given Args (and Ellipsis)
type CallExpr struct {
	Fun      Expr
	Args     []Expr // function calling arguments
	Ellipsis bool   // in go, there may be an ellipsis declaration to "unpack" a slice of values into varargs.
	Obj
}

func NewCallExpr(fun Expr, args ...Expr) *CallExpr {
	n := &CallExpr{
		Fun:  fun,
		Args: args,
	}

	assertNotAttached(fun)
	assertSettableParent(fun).SetParent(n)

	for _, arg := range args {
		assertNotAttached(arg)
		assertSettableParent(arg).SetParent(n)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *CallExpr) Children() []Node {
	tmp := make([]Node, 0, len(n.Args)+1)
	tmp = append(tmp, n.Fun)
	for _, arg := range n.Args {
		tmp = append(tmp, arg)
	}

	return tmp
}

func (n *CallExpr) exprNode() {

}
