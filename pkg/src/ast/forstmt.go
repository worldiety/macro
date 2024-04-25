package ast

// ForStmt declares a for loop.
type ForStmt struct {
	Init Node // actually a statement, like a variable definition. May be nil
	Cond Node // condition or nil
	Post Node // post iteration or nil
	Body *Block
	Obj
}

func NewForStmt(init, cond, post Node, body *Block) *ForStmt {
	n := &ForStmt{}

	n.Init = init
	if init != nil {
		assertNotAttached(init)
		assertSettableParent(init).SetParent(n)
	}

	n.Cond = cond
	if cond != nil {
		assertNotAttached(cond)
		assertSettableParent(cond).SetParent(n)
	}

	n.Post = post
	if post != nil {
		assertNotAttached(post)
		assertSettableParent(post).SetParent(n)
	}

	n.Body = body
	if body != nil {
		assertNotAttached(body)
		assertSettableParent(body).SetParent(n)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *ForStmt) Children() []Node {
	tmp := make([]Node, 0, 4)
	if n.Init != nil {
		tmp = append(tmp, n.Init)
	}

	if n.Cond != nil {
		tmp = append(tmp, n.Init)
	}

	if n.Post != nil {
		tmp = append(tmp, n.Init)
	}

	if n.Body != nil {
		tmp = append(tmp, n.Body)
	}

	return tmp
}
