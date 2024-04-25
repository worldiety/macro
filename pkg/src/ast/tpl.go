package ast

// Tpl contains uninterpreted code and should be used with enormous care. It is passed on an as-is basis to
// the renderer and may break severely. It cannot be inspected, transformed or validated in any way. Especially
// automatic imports (and thus avoiding collisions) will not be done.
//
// Why is this useful? There may be large templated parts like helper code which is complex and inspecting
// or transpiling the implementation is not relevant or possible.
//
// The given Template is applied with a go text template processor, to allow access to some render helpers, like
// the importer etc. See also TplContext.
type Tpl struct {
	Template string
	Values   map[string]interface{}
	Obj
}

func NewTpl(tpl string) *Tpl {
	return &Tpl{Template: tpl}
}

// Put inserts a key-value pair which can be accessed later by TplContext.
func (n *Tpl) Put(key string, val interface{}) *Tpl {
	if n.Values == nil {
		n.Values = map[string]interface{}{}
	}

	n.Values[key] = val

	return n
}

func (n *Tpl) exprNode() {

}

// TplContext is provided to the Template and can be used within the go template syntax.
// For example a declaration like {{ .Use "context.Context" }} will cause an import.
// Go:
//   fmt.Println({{ .Use "unsafe.Pointer"}}(x))
//   fmt.Println({{.Get "var"}})
type TplContext interface {
	// Use interprets the given name as an ast.Name and marks it for the importer. It may rewrite and return
	// another ast.Name which must be used instead.
	Use(name string) string

	// Self returns the Tpl reference.
	Self() *Tpl

	// Get returns the uninterpreted value of Tpl.Values or nil if undefined.
	Get(key string) interface{}
}
