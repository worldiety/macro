package ast

import "reflect"

// File represents a physical source code file respective compilation unit.
//  * Go: <lowercase AnnotationName>.go
//  * Java: <CamelCasePrimaryTypeName>.java
type File struct {
	// A Preamble comment belongs not to any type and is usually
	// something like a license or generator header as the first comment In the actual file.
	// The files comment is actually Obj.Comment.
	Preamble *Comment
	Name     string
	Nodes    []Node
	Obj
}

// NewFile allocates a new File.
func NewFile(name string) *File {
	return &File{Name: name}
}

// SetPreamble sets a non-package comment.
func (n *File) SetPreamble(text string) *File {
	n.Preamble = NewComment(text)
	n.Preamble.SetParent(n)
	return n
}

// SetComment sets the package comment section.
func (n *File) SetComment(text string) *File {
	n.ObjComment = NewComment(text)
	n.ObjComment.SetParent(n)
	return n
}

func (n *File) AddTypes(t ...Node) *File {
	for _, node := range t {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(n)
		n.Nodes = append(n.Nodes, node)
	}

	return n
}

// Types returns all kinds of NamedType (e.g. Enum, Struct or Interface).
func (n *File) Types() []NamedType {
	var res []NamedType

	for _, node := range n.Nodes {
		if f, ok := node.(NamedType); ok {
			res = append(res, f)
		}
	}

	return res
}

// Pkg asserts that the parent is a Pkg instance and returns it.
func (n *File) Pkg() *Pkg {
	if p, ok := n.Parent().(*Pkg); ok {
		return p
	}

	panic("expected parent to be a *Pkg, but was: " + reflect.TypeOf(n.Parent()).Name())
}

func (n *File) AddFuncs(t ...*Func) *File {
	for _, node := range t {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(n)
		n.Nodes = append(n.Nodes, node)
	}

	return n
}

func (n *File) AddNodes(t ...Node) *File {
	for _, node := range t {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(n)
		n.Nodes = append(n.Nodes, node)
	}

	return n
}

// Funcs returns all declared functions.
func (n *File) Funcs() []*Func {
	var res []*Func

	for _, node := range n.Nodes {
		if f, ok := node.(*Func); ok {
			res = append(res, f)
		}
	}

	return res
}

// Interfaces returns all declared interfaces.
func (n *File) Interfaces() []*Interface {
	var res []*Interface

	for _, node := range n.Nodes {
		if f, ok := node.(*Interface); ok {
			res = append(res, f)
		}
	}

	return res
}

// Imports returns all manually declared imports. Automatically derived imports are not returned.
func (n *File) Imports() []*Import {
	var res []*Import

	for _, node := range n.Nodes {
		if f, ok := node.(*Import); ok {
			res = append(res, f)
		}
	}

	return res
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *File) Children() []Node {
	return append([]Node{}, n.Nodes...)
}
