package ast

import (
	"reflect"
)

var nodeType = reflect.TypeOf((*Node)(nil)).Elem()

// ParentAs starts at the given node and walks up the parent hierarchy until the first found node is assignable to
// target or no more parents exists. Example:
//   mod := &ast.Mod{}
//   if ok := ast.ParentAs(someNode, &mod); ok{
//   ...
//   }
func ParentAs(node Node, target interface{}) bool {
	if target == nil {
		panic("ast: target cannot be nil")
	}

	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("ast: target must be a non-nil pointer")
	}

	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(nodeType) {
		panic("ast: *target must be interface or implement yast.Node")
	}

	targetType := typ.Elem()
	for node != nil {
		if reflect.TypeOf(node).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(node))
			return true
		}

		node = node.Parent()
	}

	return false
}

// Root returns the top most parent node or self if already root.
func Root(node Node) Node {
	if node.Parent() == nil {
		return node
	}

	return Root(node.Parent())
}

// ForEachMod loops until the parent and descends to loop all available modules.
func ForEachMod(node Node, f func(mod *Mod) error) error {
	root := Root(node)
	switch t := root.(type) {
	case *Prj:
		for _, module := range t.Mods {
			if err := f(module); err != nil {
				return err
			}
		}
	case *Mod:
		return f(t)
	}

	return nil
}

func assertNotAttached(n Node) {
	if n.Parent() != nil {
		panic("assert: node " + reflect.TypeOf(n).String() + " is already attached to " + reflect.TypeOf(n.Parent()).String())
	}
}

func assertSettableParent(node Node) SettableParent {
	if sp, ok := node.(SettableParent); ok {
		return sp
	} else {
		panic("assert: node must be a SettableParent: " + reflect.TypeOf(node).String())
	}
}

// ForEach walks recursively over each node and children.
func ForEach(parent Node, f func(n Node) error) error {
	if err := f(parent); err != nil {
		return err
	}

	if p, ok := parent.(Parent); ok {
		for _, node := range p.Children() {
			if err := ForEach(node, f); err != nil {
				return err
			}
		}
	}

	return nil
}

