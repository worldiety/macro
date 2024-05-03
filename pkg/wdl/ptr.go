package wdl

type MutPtr struct {
	pkg     *Package
	tDef    TypeDef
	macros  []*MacroInvocation
	comment []*CommentLine
}

func (m *MutPtr) Pkg() *Package {
	return m.pkg
}

func (m *MutPtr) SetPkg(pkg *Package) {
	m.pkg = pkg
}

func (m *MutPtr) SetMacros(macros []*MacroInvocation) {
	m.macros = macros
}

func (m *MutPtr) Comment() []*CommentLine {
	return m.comment
}

func (m *MutPtr) SetComment(comment []*CommentLine) {
	m.comment = comment
}

func (m *MutPtr) typeDef() {
}

func (m *MutPtr) Name() Identifier {
	return ""
}

func (m *MutPtr) AsResolvedType() *ResolvedType {
	rtype := m.tDef.AsResolvedType()
	rtype.SetPointer(true)
	return rtype
}

func (m *MutPtr) Macros() []*MacroInvocation {
	return m.macros
}

func (m *MutPtr) TypeDef() TypeDef {
	return m.tDef
}

func (m *MutPtr) SetTypeDef(typeDef TypeDef) {
	m.tDef = typeDef
}

func NewMutPtr(with func(ptr *MutPtr)) *MutPtr {
	p := &MutPtr{}
	if with != nil {
		with(p)
	}
	return p
}
