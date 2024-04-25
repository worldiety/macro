package ast

// A Field represents a (usually named) attribute or member of a struct or class.
type Field struct {
	FieldVisibility  Visibility
	FieldName        string
	FieldType        TypeDecl
	FieldAnnotations []*Annotation
	FieldDefault     *BasicLit
	Obj
}

// NewField allocates a new named field. For some renderers like Go, an empty name declares an embedded type.
func NewField(name string, typeDecl TypeDecl) *Field {
	f := &Field{
		FieldName: name,
		FieldType: typeDecl,
	}

	assertNotAttached(typeDecl)
	assertSettableParent(typeDecl).SetParent(f)

	return f
}

// SetDefault attaches the literal to be a field value initializer. This is not supported by all languages (e.g.
// Go does not have that).
func (f *Field) SetDefault(lit *BasicLit) *Field {
	assertNotAttached(lit)
	assertSettableParent(lit).SetParent(f)

	if f.FieldDefault != nil {
		f.FieldDefault.SetParent(nil)
	}

	f.FieldDefault = lit

	return f
}

// SetComment sets the nodes comment.
func (f *Field) SetComment(text string) *Field {
	f.ObjComment = NewComment(text)
	f.ObjComment.SetParent(f)
	return f
}

// SetVisibility updates the fields Visibility. The Go renderer will override the rendered name to match the visibility.
func (f *Field) SetVisibility(v Visibility) *Field {
	f.FieldVisibility = v
	return f
}

// Visibility returns the fields visibility.
func (f *Field) Visibility() Visibility {
	return f.FieldVisibility
}

// SetName updates the fields name.
func (f *Field) SetName(name string) *Field {
	f.FieldName = name
	return f
}

// Name returns the fields name.
func (f *Field) Identifier() string {
	return f.FieldName
}

// AddAnnotations appends the given annotations or tags to the field.
func (f *Field) AddAnnotations(a ...*Annotation) *Field {
	for _, annotation := range a {
		assertNotAttached(annotation)
		assertSettableParent(annotation).SetParent(f)
		f.FieldAnnotations = append(f.FieldAnnotations, annotation)
	}

	return f
}

// Annotations returns the backing slice of all annotations.
func (f *Field) Annotations() []*Annotation {
	return f.FieldAnnotations
}

// TypeDecl returns the current type declaration.
func (f *Field) TypeDecl() TypeDecl {
	return f.FieldType
}

// String returns a debugging representation.
func (f *Field) String() string {
	return f.FieldName + " " + f.FieldType.String()
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (f *Field) Children() []Node {
	tmp := make([]Node, 0, len(f.FieldAnnotations)+1+1)
	tmp = append(tmp, f.FieldType)
	for _, annotation := range f.FieldAnnotations {
		tmp = append(tmp, annotation)
	}

	if f.FieldDefault != nil {
		tmp = append(tmp, f.FieldDefault)
	}

	return tmp
}
