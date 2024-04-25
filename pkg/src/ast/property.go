package ast

// A Property represents a (usually named) attribute or member of a struct or class. The field has always
// private semantics.
//
// Go
//  renders as a private field with public getter.
// TODO is this actually required?
type Property struct {
	FieldName string
	FieldType TypeDecl
	Obj

	Read struct {
		Enabled    bool
		Visibility Visibility
	}

	Write struct {
		Enabled    bool
		Visibility Visibility
	}
}

func NewProperty(name string, fieldType TypeDecl) *Property {
	p := &Property{FieldName: name, FieldType: fieldType}
	assertNotAttached(fieldType)
	assertSettableParent(fieldType).SetParent(p)

	return p
}

// SetComment sets the nodes comment.
func (p *Property) SetComment(text string) *Property {
	p.ObjComment = NewComment(text)
	p.ObjComment.SetParent(p)
	return p
}

func (p *Property) Reader(enabled bool, visibility Visibility) *Property {
	p.Read.Enabled = enabled
	p.Read.Visibility = visibility

	return p
}

func (p *Property) Writer(enabled bool, visibility Visibility) *Property {
	p.Write.Enabled = enabled
	p.Write.Visibility = visibility

	return p
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (p *Property) Children() []Node {
	return []Node{p.FieldType}
}
