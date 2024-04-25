package ast

// An Expr is a specialised node. Due to our macros, this is expressiveness is useless and due to different
// languages wrong anyway, like assignments (which are statements in Go and expressions in Java).
type Expr interface {
	Node
	exprNode() // marker interface method
}

// Exprs is a builder helper for varargs. It filters out nil expressions
func Exprs(expr ...Expr) []Expr {
	tmp := make([]Expr, 0, len(expr))
	for _, e := range expr {
		if e != nil{
			tmp = append(tmp, e)
		}
	}

	return tmp
}
