package ast

// A positioned Comment text. The Text may have been normalized and prefixes and others have been removed, so
// an exact positioning of the actual text characters is not possible. A comment may consists of multiple
// single line comments or one large multiline comment.
//   * Go/Java: // for single line and /* .. */ for multiline, however the text is stripped.
type Comment struct {
	Text string // the actual comment text, may include newlines.
	Obj
}

// NewComment allocates a new Comment node.
func NewComment(text string) *Comment {
	return &Comment{Text: text}
}

func (n *Comment) Clone() *Comment {
	if n == nil {
		return nil
	}

	c := &Comment{
		Text: n.Text,
	}

	if n.ObjComment != nil {
		c.ObjComment = n.ObjComment.Clone()
	}

	return c
}
