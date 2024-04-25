package ast

// IfStmt describes something like
//  Go: if x > y {
//      }
//
//      or
//
//      if a, err := d.call(); err!=nil {
//      }
type IfStmt struct {
	Init Node // actually a statement, like a variable definition. May be nil
	Cond Expr
	Body *Block
	Else Node // actually a statement. May be nil.
	Obj
}

func NewIfStmt(cond Expr, body *Block) *IfStmt {
	n := &IfStmt{
		Cond: cond,
		Body: body,
	}

	assertNotAttached(cond)
	assertSettableParent(cond).SetParent(n)

	assertNotAttached(body)
	assertSettableParent(body).SetParent(n)

	return n
}

func (n *IfStmt) SetInit(initStmt Node) *IfStmt {
	n.Init = initStmt
	assertNotAttached(initStmt)
	assertSettableParent(initStmt).SetParent(n)

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *IfStmt) Children() []Node {
	tmp := make([]Node, 0, 4)
	if n.Init != nil {
		tmp = append(tmp, n.Init)
	}

	tmp = append(tmp, n.Cond)

	if n.Body != nil {
		tmp = append(tmp, n.Body)
	}

	if n.Else != nil {
		tmp = append(tmp, n.Else)
	}

	return tmp
}
