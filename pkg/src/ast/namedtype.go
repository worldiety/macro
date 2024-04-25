package ast

// NamedType declares a new named type which is either derived from another existing type
// or is a struct or interface.
type NamedType interface {
	// Name returns the name of the declared type. Note that it may be also used (slightly incorrect) in non-named
	// situations like anonymous types by just keeping the name empty.
	Identifier() string
	Node
	sealedNamedType()
}
