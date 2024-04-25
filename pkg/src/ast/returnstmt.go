package ast

// ReturnStmt represents a single or multi value return.
type ReturnStmt struct {
	Results []Expr
	Obj
}

func NewReturnStmt(results ...Expr) *ReturnStmt {
	n := &ReturnStmt{Results: results}
	for _, result := range results {
		assertNotAttached(result)
		assertSettableParent(result).SetParent(n)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *ReturnStmt) Children() []Node {
	tmp := make([]Node, 0, len(n.Results))
	for _, result := range n.Results {
		tmp = append(tmp, result)
	}

	return tmp
}
