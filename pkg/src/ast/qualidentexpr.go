package ast

// QualifierIdent denotes an imported qualifier. See also Name.
type QualIdent struct {
	Qualifier string // Qualifier is just the qualifying part of a Name.
	Obj
}

// NewQualIdent creates a new qualified identifier. This usually refers to the package name which contains
// other identifiers. The actual renderer will decide how to reference the identifier, like a renamed
// import or full qualified usage. Examples:
//  Go: "fmt" or "http/url"
//  Java: "java.util"
func NewQualIdent(qualifier string) *QualIdent {
	return &QualIdent{Qualifier: qualifier}
}

func (n *QualIdent) exprNode() {

}
