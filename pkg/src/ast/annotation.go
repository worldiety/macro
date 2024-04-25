package ast

import "sort"

// An Annotation represents an AnnotationName and a bunch of named Values. The Go renderer will emit this as a struct field
// tag. However, only the value for the empty key is used, just as is (but quoted). In Java each key represents
// the named attribute of an annotation. The AnnotationName is interpreted as a fully qualified identifier.
type Annotation struct {
	AnnotationName Name
	Values         map[string]string
	Obj
}

// NewAnnotation creates a new named Annotation. In Go the AnnotationName is just interpreted as a string and has no further
// meaning.
func NewAnnotation(name Name) *Annotation {
	return &Annotation{
		AnnotationName: name,
		Values:         map[string]string{},
	}
}

// SetIdentifier updates the ParamAnnotations AnnotationName.
func (a *Annotation) SetIdentifier(name Name) *Annotation {
	a.AnnotationName = name
	return a
}

// Identifier returns the ParamAnnotations AnnotationName.
func (a *Annotation) Identifier() Name {
	return a.AnnotationName
}

// SetDefault sets the unnamed attribute value. See SetValue.
func (a *Annotation) SetDefault(value string) *Annotation {
	return a.PutLiteral("", value)
}

// SetValue sets a named attribute value. The value is interpreted as is, so e.g. use plain
// Values for language constants, like 3, 3.4, true or "hello world" In Java. The Go renderer only
// ever evaluates the unnamed attribute and quotes the string itself.
func (a *Annotation) PutLiteral(name, value string) *Annotation {
	a.Values[name] = value
	return a
}

// Value returns a specific value or the empty string.
func (a *Annotation) GetLiteral(name string) string {
	return a.Values[name]
}

// Attributes returns the ascended sorted list of attribute names.
func (a *Annotation) Attributes() []string {
	var tmp []string
	for key := range a.Values {
		tmp = append(tmp, key)
	}

	sort.Strings(tmp)
	return tmp
}
