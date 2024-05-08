package wdl

type Comment struct {
	lines  []*CommentLine
	macros []*MacroInvocation
}

func (c *Comment) String() string {
	if c == nil {
		return ""
	}
	tmp := ""
	for _, l := range c.lines {
		tmp += l.Text()
		tmp += "\n"
	}

	return tmp
}

func (c *Comment) Macros() []*MacroInvocation {
	if c == nil {
		return nil
	}
	return c.macros
}

func (c *Comment) SetMacros(macros []*MacroInvocation) {
	c.macros = macros
}

func (c *Comment) AddMacros(macros ...*MacroInvocation) {
	c.macros = append(c.macros, macros...)
}

func (c *Comment) Lines() []*CommentLine {
	if c == nil {
		return nil
	}
	return c.lines
}

func (c *Comment) SetLines(lines []*CommentLine) {
	c.lines = lines
}

func (c *Comment) AddLines(lines ...*CommentLine) {
	c.lines = append(c.lines, lines...)
}

func NewSimpleComment(text string) *Comment {
	return NewComment(func(comment *Comment) {
		comment.AddLines(NewCommentLine(func(line *CommentLine) {
			line.SetText(text)
		}))
	})
}

func NewComment(with func(comment *Comment)) *Comment {
	comment := &Comment{}
	if with != nil {
		with(comment)
	}

	return comment
}

type CommentLine struct {
	text string
	pos  Pos
}

func NewCommentLine(with func(line *CommentLine)) *CommentLine {
	c := &CommentLine{}
	if with != nil {
		with(c)
	}

	return c
}

func (c *CommentLine) Text() string {
	return c.text
}

func (c *CommentLine) SetText(text string) {
	c.text = text
}

func (c *CommentLine) Pos() Pos {
	return c.pos
}

func (c *CommentLine) SetPos(pos Pos) {
	c.pos = pos
}
