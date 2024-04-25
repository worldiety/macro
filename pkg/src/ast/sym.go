package ast

type SymKind int

const (
	SymTermStmt SymKind = 1 // terminator symbol for statements
	SymNewline  SymKind = 2 // newline symbol, some language also conclude a SymTermStmt from this.
)

// Sym is used to encode information which do not belong to an AST, like terminator symbols. However, I don't
// want to use the Obj positions to encode that, because our use case is usually generating code from
// scratch and not from parsing. See also https://github.com/golang/go/blob/master/src/go/printer/printer.go#L979.
type Sym struct {
	Kind SymKind
	Obj
}

func NewSym(kind SymKind) *Sym {
	return &Sym{Kind: kind}
}
