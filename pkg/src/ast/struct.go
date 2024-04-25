package ast

var _ NamedType = (*Struct)(nil)

// A Struct is actually a data type, like a record. Depending on the language, it can be used in a value or reference
// context. If supported, the primary use case should be the usage as a value to improve conclusiveness and
// performance by avoiding heap allocation (and potentially GC overhead). Inheritance is not possible, but other
// types may be embedded (e.g. in Go). Languages like Java use just simple classes (PoJos), because records have no
// exclusive use (they are just syntax sugar for a class with final members). In contrast to that, Go cannot express
// final fields.
type Struct struct {
	TypeName        string
	TypeVisibility  Visibility
	TypeFields      []*Field
	TypeStatic      bool
	TypeAnnotations []*Annotation
	TypeMethods     []*Func
	Types           []NamedType // only valid for language which can declare named nested type like java
	Implements      []Name      // Implements denotes a bunch of interfaces which must be implemented by this struct. Depending on the renderer (like Go) this has no effect.
	Embedded        []TypeDecl  // Embedded is only valid for languages which supports composition at a language level
	FactoryRefs     []*Func     // FactoryRefs are NOT considered children of a struct. They are still connected to a file, however they are considered to be a kind of constructor.
	DefaultRecName  string      // useful to transport a standard receiver name. However, you need to care yourself.
	Obj
}

// NewStruct returns a new named struct type. A struct is always mutable, but may be used either in a value
// or pointer context. Structs are straightforward in Go but in Java just a PoJo. We do not use records, because
// they have a different semantic (read only).
func NewStruct(name string) *Struct {
	return &Struct{TypeName: name}
}

func (s *Struct) AddEmbedded(t ...TypeDecl) *Struct {
	for _, decl := range t {
		assertNotAttached(decl)
		assertSettableParent(decl).SetParent(s)

		s.Embedded = append(s.Embedded, decl)
	}

	return s
}

// SetComment sets the nodes comment.
func (s *Struct) SetComment(text string) *Struct {
	s.ObjComment = NewComment(text)
	s.ObjComment.SetParent(s)
	return s
}

func (s *Struct) SetDefaultRecName(n string) *Struct {
	s.DefaultRecName = n
	return s
}

// Static returns true, if this struct or class should pull its outer scope. This is only for Java and inner classes.
func (s *Struct) Static() bool {
	return s.TypeStatic
}

// SetStatic updates the static flag. Only for Java.
func (s *Struct) SetStatic(static bool) *Struct {
	s.TypeStatic = static
	return s
}

// Identifier returns the declared identifier which must be unique per package.
func (s *Struct) Identifier() string {
	return s.TypeName
}

func (s *Struct) sealedNamedType() {
	panic("implement me")
}

// SetVisibility sets the visibility. The default is Public.
func (s *Struct) SetVisibility(v Visibility) *Struct {
	s.TypeVisibility = v
	return s
}

// Visibility returns the current visibility. The default is Public.
func (s *Struct) Visibility() Visibility {
	return s.TypeVisibility
}

// AddFields appends the given fields to the struct.
func (s *Struct) AddFields(fields ...*Field) *Struct {
	for _, field := range fields {
		assertNotAttached(field)
		assertSettableParent(field).SetParent(s)
		s.TypeFields = append(s.TypeFields, field)
	}
	return s
}

// AddFactoryRefs just appends the given funcs for the purpose of factories or constructors. Most importantly
// Struct does not take the ownership and the parent is still unset (usually a file or another type).
func (s *Struct) AddFactoryRefs(f ...*Func) *Struct {
	s.FactoryRefs = append(s.FactoryRefs, f...)

	return s
}

// Fields returns the currently configured fields.
func (s *Struct) Fields() []*Field {
	return s.TypeFields
}

// Annotations returns the backing slice of all annotations.
func (s *Struct) Annotations() []*Annotation {
	return s.TypeAnnotations
}

// AddAnnotations appends the given annotations. Note that not all render targets support type annotations, e.g.
// like Go.
func (s *Struct) AddAnnotations(a ...*Annotation) *Struct {
	for _, annotation := range a {
		assertNotAttached(annotation)
		assertSettableParent(annotation).SetParent(s)
	}

	return s
}

// Methods returns all available functions.
func (s *Struct) Methods() []*Func {
	return s.TypeMethods
}

// AddMethods appends more methods to this interfaces contract.
func (s *Struct) AddMethods(f ...*Func) *Struct {
	for _, fun := range f {
		assertNotAttached(fun)
		assertSettableParent(fun).SetParent(s)
		s.TypeMethods = append(s.TypeMethods, fun)
	}

	return s
}

func (s *Struct) NamedTypes() []NamedType {
	return s.Types
}

func (s *Struct) AddNamedTypes(t ...NamedType) *Struct {
	for _, namedType := range t {
		assertNotAttached(namedType)
		assertSettableParent(namedType).SetParent(s)
		s.Types = append(s.Types, namedType)
	}

	return s
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
// FactoryRefs are not considered children, to avoid recursive loops in the AST.
func (s *Struct) Children() []Node {
	tmp := make([]Node, 0, len(s.TypeFields)+len(s.TypeAnnotations)+len(s.TypeMethods)+len(s.Types)+len(s.Embedded))
	for _, param := range s.TypeAnnotations {
		tmp = append(tmp, param)
	}

	for _, param := range s.TypeFields {
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
