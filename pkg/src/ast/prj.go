package ast

// A Prj usually refers to a kind of workspace or to the root of a mono repo. In distributed systems,
// this may be entirely artificially and has no physical correspondence (beside a potential architecture project).
// Most importantly this is the root of heterogen modules, like a go server, a java library and a swift application.
type Prj struct {
	Name string
	Mods []*Mod

	Obj
}

// NewPrj allocates a new project.
func NewPrj(name string) *Prj {
	return &Prj{Name: name}
}

// FindProject walks up the node hierarchy until it finds a project to return. Otherwise returns nil.
func FindProject(n Node) *Prj {
	prj := &Prj{}
	if ok := ParentAs(n, &prj); ok {
		return prj
	}

	return nil
}

// AddModules appends and attaches the given modules.
func (n *Prj) AddModules(m ...*Mod) *Prj {
	for _, mod := range m {
		assertNotAttached(mod)
		assertSettableParent(mod).SetParent(n)
		n.Mods = append(n.Mods, mod)
	}

	return n
}
