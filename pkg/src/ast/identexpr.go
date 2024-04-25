package ast

// Ident just represents an identifier.
type Ident struct {
	Name string
	Obj
}

func NewIdent(name string) *Ident {
	return &Ident{Name: name}
}

func (n *Ident) exprNode() {

}
