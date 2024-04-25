package ast

import "strconv"

// A Pos describes a resolved position within a file.
type Pos struct {
	// File contains the absolute file path.
	File string
	// Line denotes the one-based line number In the denoted File.
	Line int
	// Col denotes the one-based column number In the denoted Line.
	Col int
}

// String returns the content In the "file:line:col" format.
func (p Pos) String() string {
	return p.File + ":" + strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Col)
}

// A Node represents the common contract
type Node interface {
	// Pos returns the actual starting position of this Node.
	Pos() Pos

	// End is the position of the first char after the node.
	End() Pos

	// Parent returns the parent Node or nil if undefined. This recursive implementation may be considered as
	// unnecessary and even as an anti pattern within an AST but the core feature is to perform semantic validations
	// which requires a lot of down/up iterations through the (entire) AST. Keeping the relational relation
	// at the node level keeps things simple and we don't need to pass (path) contexts everywhere.
	Parent() Node

	// Value is like a context Value getter.
	Value(key interface{}) interface{}

	// PutValue overwrites the given key with the given value. Consider using package private types, to ensure
	// that no collisions can occur.
	PutValue(key, value interface{})

	// Comment returns an optional comment node.
	Comment() *Comment
}

func Nodes(n ...Node) []Node {
	return n
}

// A Parent is a Node and may contain other nodes as children. This is used to simplify algorithms based on Walk.
type Parent interface {
	Node
	// Children returns a defensive copy of the underlying slice. However the Node references are shared.
	Children() []Node
}

// Obj is actually a helper to implement a Node by embedding the Obj
type Obj struct {
	ObjPos     Pos
	ObjEnd     Pos
	ObjParent  Node
	ObjComment *Comment // the actual comment of the logical object
	Values     map[interface{}]interface{}
}

func (n *Obj) Pos() Pos {
	return n.ObjPos
}

func (n *Obj) End() Pos {
	return n.ObjEnd
}

func (n *Obj) Parent() Node {
	return n.ObjParent
}

func (n *Obj) SetParent(p Node) {
	n.ObjParent = p
}

func (n *Obj) Comment() *Comment {
	return n.ObjComment
}

// CommentText returns either the empty string or the comments text.
func (n *Obj) CommentText() string {
	if n == nil || n.ObjComment == nil {
		return ""
	}

	return n.ObjComment.Text
}

// Value is like a context Value getter.
func (n *Obj) Value(key interface{}) interface{} {
	if n.Values == nil {
		return nil
	}

	return n.Values[key]
}

// PutValue overwrites the given key with the given value. Consider using package private types, to ensure
// that no collisions can occur.
func (n *Obj) PutValue(key, value interface{}) {
	if n.Values == nil {
		n.Values = map[interface{}]interface{}{}
	}

	n.Values[key] = value

}

func (n *Obj) Clone() *Obj {
	values := make(map[interface{}]interface{})
	for k, v := range n.Values {
		values[k] = v
	}

	return &Obj{
		ObjPos:     n.ObjPos,
		ObjEnd:     n.ObjEnd,
		ObjParent:  nil,
		ObjComment: n.ObjComment.Clone(),
		Values:     values,
	}
}

type SettableParent interface {
	Node
	SetParent(p Node)
}
