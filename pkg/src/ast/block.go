package ast

// A Block represents a lexical group of declarations and statements. Usually a block also introduces a scope and
// can also be nested.
//  Go/Java: { ... }
type Block struct {
	Nodes []Node
	Obj
}

// NewBlock allocates a new block
func NewBlock(nodes ...Node) *Block {
	n := &Block{}
	for _, node := range nodes {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(n)
		n.Nodes = append(n.Nodes, node)
	}

	return n
}

// SetComment sets the nodes comment.
func (n *Block) SetComment(text string) *Block {
	n.ObjComment = NewComment(text)
	n.ObjComment.SetParent(n)
	return n
}

// Add appends and attaches the given nodes to this block.
func (n *Block) Add(nodes ...Node) *Block {
	for _, node := range nodes {
		n.Nodes = append(n.Nodes, node)
		assertNotAttached(node)
		assertSettableParent(node).SetParent(n)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Block) Children() []Node {
	return append(make([]Node, 0, len(n.Nodes)), n.Nodes...)
}
