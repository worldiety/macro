package ast

// A CompLit represents a composite literal
type CompLit struct {
	Type     Expr   // the literal type. There may be situations where this is nil (anonymous struct literals in Go)
	Elements []Expr // list of the composite elements, may be named or not.
	Obj
}

func NewCompLit(typ Expr, elems ...Expr) *CompLit {
	n := &CompLit{
		Type:     typ,
		Elements: elems,
	}

	assertNotAttached(typ)
	assertSettableParent(typ).SetParent(n)

	n.AddElements(elems...)

	return n
}

func (n *CompLit) AddElements(elems ...Expr) *CompLit {
	for _, elem := range elems {
		assertNotAttached(elem)
		assertSettableParent(elem).SetParent(n)
		n.Elements = append(n.Elements, elem)
	}

	return n
}

func (n *CompLit) exprNode() {

}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *CompLit) Children() []Node {
	tmp := make([]Node, 0, len(n.Elements)+1)
	if n.Type != nil {
		tmp = append(tmp, n.Type)
	}

	for _, element := range n.Elements {
		tmp = append(tmp, element)
	}

	return tmp
}
