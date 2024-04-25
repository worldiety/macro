package ast

// A Macro provides a dynamic amount of children and does not cause emitted source itself. It can be used
// to emit generator specific nodes at calling time based on the current ast state (especially its parents).
type Macro struct {
	ID   string                // an arbitrary ID
	Func func(m *Macro) []Node // Func should always return a new defensive copy with shared Node instances.
	// if true, the result of Func is ever evaluated once. This improves performance but also makes stateful macros
	// easier to implement when called multiple times. However, when rendering for multiple platforms, this may
	// cause wrong results. See also Invalidate.
	CacheFunc bool
	funcCache []Node
	Obj
}

func NewMacro() *Macro {
	return &Macro{CacheFunc: true}
}

func (n *Macro) exprNode() {

}

func (n *Macro) SetID(id string) *Macro {
	n.ID = id

	return n
}

// Invalidate purges any node cache.
func (n *Macro) Invalidate() *Macro {
	n.funcCache = nil
	return n
}

// Target returns the available module target information. If there is no target available, all fields are default.
func (n *Macro) Target() Target {
	mod := &Mod{}
	if ok := ParentAs(n, &mod); ok {
		return mod.Target
	}

	return Target{}
}

// Children just delegates to Func.
func (n *Macro) Children() []Node {
	if n.Func == nil {
		return nil
	}

	if !n.CacheFunc {
		return n.Func(n)
	}

	if n.funcCache != nil {
		return n.funcCache
	}

	n.funcCache = n.Func(n)

	return n.funcCache
}

// SetMatchers is a builder function which replaces the Func with a loop implementation which invokes each given
// matcher in the given order and therefore just returns the first static nodes which apply.
func (n *Macro) SetMatchers(matchers ...func(m *Macro) (bool, []Node)) *Macro {
	n.Func = func(m *Macro) []Node {
		for _, f := range matchers {
			matches, nodes := f(n)
			if matches {
				return nodes
			}
		}

		return nil
	}

	return n
}

// MatchTargetLanguage returns a closure which can be used in conjunction with Macro.SetMatchers and
// evaluates to true as soon as the target language matches the given language. The given static nodes
// are just returned. Note that each node is attached to this macro on successful evaluation.
func MatchTargetLanguage(lang Lang, nodes ...Node) func(m *Macro) (bool, []Node) {
	return MatchTargetLanguageWithContext(lang, func(m *Macro) []Node {
		return nodes
	})
}

func MatchTargetLanguageWithContext(lang Lang, f func(m *Macro) []Node) func(m *Macro) (bool, []Node) {
	return func(m *Macro) (bool, []Node) {
		nodes := f(m)

		target := m.Target()
		if target.Lang == lang {
			for _, node := range nodes {
				if node.Parent() != nil && node.Parent() != m {
					assertNotAttached(node)
				}

				if node.Parent() == nil {
					assertSettableParent(node).SetParent(m)
				}

			}
			return true, nodes
		}

		return false, nil
	}
}
