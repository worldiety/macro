package wdl

// A TypeDef is a marker interface and contains all types which define the memory layout of a type.
// Note, that this is the c semantic and not the Go Spec semantic, which uses the term declaration for functions and
// other GenDecl like alias, interface et al. They use the term type definition when introducing a new defined
// type based on another type (e.g. type Text string)
type TypeDef interface {
	typeDef()
	Name() Identifier
	AsResolvedType() *ResolvedType
	Macros() []*MacroInvocation
	Clone() TypeDef
	SetTypeParams(typeParams []*ResolvedType)
}
