package ast

// RangeStmt declares a for loop.
type RangeStmt struct {
	Key   Node // may be nil
	Val Node // may be nil
	X     Node // range target
	Body  *Block
	Obj
}

func NewRangeStmt(key, val, x Node, body *Block) *RangeStmt {
	n := &RangeStmt{}

	n.Key = key
	if key != nil {
		assertNotAttached(key)
		assertSettableParent(key).SetParent(n)
	}

	n.Val = val
	if val != nil {
		assertNotAttached(val)
		assertSettableParent(val).SetParent(n)
	}

	n.X = x
	if x != nil {
		assertNotAttached(x)
		assertSettableParent(x).SetParent(n)
	}

	n.Body = body
	if body != nil {
		assertNotAttached(body)
		assertSettableParent(body).SetParent(n)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *RangeStmt) Children() []Node {
	tmp := make([]Node, 0, 4)
	if n.Key != nil {
		tmp = append(tmp, n.Key)
	}

	if n.Val != nil {
		tmp = append(tmp, n.Val)
	}

	if n.X != nil {
		tmp = append(tmp, n.X)
	}

	if n.Body != nil {
		tmp = append(tmp, n.Body)
	}

	return tmp
}
