package wdl

type Interface struct {
	pkg     *Package
	name    Identifier
	macros  []*MacroInvocation
	comment []*CommentLine
	types   []*ResolvedType
}

func (u *Interface) Pkg() *Package {
	return u.pkg
}

func (u *Interface) SetPkg(pkg *Package) {
	u.pkg = pkg
}

func (u *Interface) AsResolvedType() *ResolvedType {
	return NewResolvedType(func(rType *ResolvedType) {
		rType.SetName(u.name)
		rType.SetPkg(u.pkg)
		rType.SetTypeDef(u)
	})
}

func (u *Interface) Name() Identifier {
	return u.name
}

func (u *Interface) SetName(name Identifier) {
	u.name = name
}

func (u *Interface) Types() []*ResolvedType {
	return u.types
}

func (u *Interface) AddTypes(types ...*ResolvedType) {
	u.types = append(u.types, types...)
}

func (u *Interface) SetTypes(types []*ResolvedType) {
	u.types = types
}

func NewInterface(with func(iface *Interface)) *Interface {
	u := &Interface{}
	if with != nil {
		with(u)
	}
	return u
}

func (u *Interface) typeDef() {}

func (u *Interface) Macros() []*MacroInvocation {
	return u.macros
}

func (u *Interface) AddMacro(macro *MacroInvocation) {
	u.macros = append(u.macros, macro)
}

func (u *Interface) SetMacros(macro []*MacroInvocation) {
	u.macros = macro
}

func (u *Interface) Comment() []*CommentLine {
	return u.comment
}

func (u *Interface) SetComment(comment []*CommentLine) {
	u.comment = comment
}
