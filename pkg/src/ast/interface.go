package ast

var _ NamedType = (*Interface)(nil)

// An Interface is a contract which defines a method set and allows polymorphism without inheritance. If a
// body is declared, it depends on the actual renderer, if a default method will be emitted (e.g. for Java).
// Go does not support default methods in interfaces.
type Interface struct {
	TypeName        string
	TypeVisibility  Visibility
	TypeMethods     []*Func
	TypeAnnotations []*Annotation
	Types           []NamedType // only valid for language which can declare named nested type like java
	Embedded        []TypeDecl  // Embedded is only valid for languages which supports composition at a language level
	Obj
}

// NewInterface returns a new named struct type. A struct is always mutable, but may be used either in a value
// or pointer context. Structs are straightforward in Go but in Java just a PoJo. We do not use records, because
// they have a different semantic (read only).
func NewInterface(name string) *Interface {
	return &Interface{TypeName: name}
}

// SetComment sets the nodes comment.
func (s *Interface) SetComment(text string) *Interface {
	s.ObjComment = NewComment(text)
	s.ObjComment.SetParent(s)
	return s
}

// Identifier returns the declared identifier which must be unique per package.
func (s *Interface) Identifier() string {
	return s.TypeName
}

// SetName updates the interfaces identifier which must be unique per package.
func (s *Interface) SetName(name string) *Interface {
	s.TypeName = name
	return s
}

func (s *Interface) sealedNamedType() {
	panic("implement me")
}

func (s *Interface) AddEmbedded(t ...TypeDecl) *Interface {
	for _, decl := range t {
		assertNotAttached(decl)
		assertSettableParent(decl).SetParent(s)

		s.Embedded = append(s.Embedded, decl)
	}

	return s
}

// SetVisibility sets the visibility. The default is Public.
func (s *Interface) SetVisibility(v Visibility) *Interface {
	s.TypeVisibility = v
	return s
}

// Visibility returns the current visibility. The default is Public.
func (s *Interface) Visibility() Visibility {
	return s.TypeVisibility
}

// Methods returns all available functions.
func (s *Interface) Methods() []*Func {
	return s.TypeMethods
}

// AddMethods appends more methods to this interfaces contract.
func (s *Interface) AddMethods(f ...*Func) *Interface {
	for _, fun := range f {
		assertNotAttached(fun)
		assertSettableParent(fun).SetParent(s)
		s.TypeMethods = append(s.TypeMethods, fun)
	}

	return s
}

// Annotations returns the backing slice of all annotations.
func (s *Interface) Annotations() []*Annotation {
	return s.TypeAnnotations
}

// AddAnnotations appends the given annotations. Note that not all render targets support type annotations, e.g.
// like Go.
func (s *Interface) AddAnnotations(a ...*Annotation) *Interface {
	for _, annotation := range a {
		assertNotAttached(annotation)
		assertSettableParent(annotation).SetParent(s)
		s.TypeAnnotations = append(s.TypeAnnotations, annotation)
	}

	return s
}

// AddNamedTypes adds a bunch of named types. This is only allowed in Java and other renderers should
// either ignore it or place them at the package level (Go).
func (s *Interface) AddNamedTypes(types ...NamedType) *Interface {
	for _, namedType := range types {
		assertNotAttached(namedType)
		assertSettableParent(namedType).SetParent(s)
		s.Types = append(s.Types, namedType)
	}

	return s
}

// NamedTypes returns the backing slice of defined named types.
func (s *Interface) NamedTypes() []NamedType {
	return s.Types
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (s *Interface) Children() []Node {
	tmp := make([]Node, 0, +len(s.TypeAnnotations)+len(s.TypeMethods)+len(s.Types)+len(s.Embedded))
	for _, param := range s.TypeAnnotations {
		tmp = append(tmp, param)
	}

	for _, param := range s.TypeMethods {
		tmp = append(tmp, param)
	}

	for _, namedType := range s.Types {
		tmp = append(tmp, namedType)
	}

	for _, e := range s.Embedded {
		tmp = append(tmp, e)
	}

	return tmp
}
