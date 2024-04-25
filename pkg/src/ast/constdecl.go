package ast

// A ConstDecl represents an immutable compile time constant.
//  Go:
//    const X = "abc"
//    const (
//       SomeStuff MyEnum = iota
//    )
//
//  Java:
//    public static final String X = "abc
type ConstDecl struct {
	Assignments []*Assign
	Obj
}

func NewConstDecl(assigns ...*Assign) *ConstDecl {
	n := &ConstDecl{}
	n.Add(assigns...)
	return n
}

func (n *ConstDecl) Add(assigns ...*Assign) *ConstDecl {
	for _, assign := range assigns {
		assertNotAttached(assign)
		assertSettableParent(assign).SetParent(n)
		n.Assignments = append(n.Assignments, assign)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *ConstDecl) Children() []Node {
	tmp := make([]Node, 0, len(n.Assignments))
	for _, param := range n.Assignments {
		tmp = append(tmp, param)
	}

	return tmp
}
