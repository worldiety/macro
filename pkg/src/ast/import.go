package ast

// Import declares an explicit import statement. Note that the imports are usually generated automatically, but
// one can do it by hand also.
type Import struct {
	Ident string
	Name  Name
	Obj
}

func NewImport(ident string, name Name) *Import {
	return &Import{Ident: ident, Name: name}
}

// SetComment sets the nodes comment.
func (s *Import) SetComment(text string) *Import {
	s.ObjComment = NewComment(text)
	s.ObjComment.SetParent(s)
	return s
}
