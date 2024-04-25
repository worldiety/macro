package ast

// A VarDecl represents a variable statement.
//  Go:
//    var X = "abc"
//    var (
//       SomeStuff MyType
//    )
//
//  Java:
//    public String X = "abc
type VarDecl struct {
	Decl []Node
	Obj
}

func NewVarDecl(assigns ...Node) *VarDecl {
	n := &VarDecl{}
	n.Add(assigns...)
	return n
}

func (n *VarDecl) Add(assigns ...Node) *VarDecl {
	for _, assign := range assigns {
		assertNotAttached(assign)
		assertSettableParent(assign).SetParent(n)
		n.Decl = append(n.Decl, assign)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *VarDecl) Children() []Node {
	tmp := make([]Node, 0, len(n.Decl))
	for _, param := range n.Decl {
		tmp = append(tmp, param)
	}

	return tmp
}
