package ast

// An Union is an algebraic sum type, whose tagged variants are enumerable.
type Union struct {
	typeName   string // TypeName denotes the actual name of this type.
	implements []Name // Implements denotes a bunch of interfaces which must be implemented by this Union. Depending on the renderer (like a specific Go version) this may break.
	types      []TypeDecl
	Obj
}

func (n *Union) Identifier() string {
	return n.typeName
}

func (n *Union) sealedNamedType() {
	panic("implement me")
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Union) Children() []Node {
	tmp := make([]Node, 0, len(n.types))

	for _, enumCase := range n.types {
		tmp = append(tmp, enumCase)
	}

	return tmp
}
