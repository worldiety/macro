package wdl

type Union struct {
	pkg        *Package
	macros     []*MacroInvocation
	comment    []*CommentLine
	types      []*ResolvedType
	name       Identifier
	file       *File
	typeParams []*ResolvedType
}

func (u *Union) Clone() TypeDef {
	return &Union{
		pkg:        u.pkg,
		macros:     append([]*MacroInvocation{}, u.macros...),
		comment:    append([]*CommentLine{}, u.comment...),
		types:      append([]*ResolvedType{}, u.typeParams...),
		name:       u.name,
		file:       u.file,
		typeParams: append([]*ResolvedType{}, u.typeParams...),
	}
}

func (u *Union) SetTypeParams(typeParams []*ResolvedType) {
	u.typeParams = typeParams
}

func (u *Union) File() *File {
	return u.file
}

func (u *Union) SetFile(file *File) {
	u.file = file
}

func (u *Union) Pkg() *Package {
	return u.pkg
}

func (u *Union) SetPkg(pkg *Package) {
	u.pkg = pkg
}

func (u *Union) AsResolvedType() *ResolvedType {
	return NewResolvedType(func(rType *ResolvedType) {
		rType.SetName(u.name)
		rType.SetPkg(u.pkg)
		rType.SetTypeDef(u)
	})
}

func (u *Union) Name() Identifier {
	return u.name
}

func (u *Union) SetName(name Identifier) {
	u.name = name
}

func (u *Union) Types() []*ResolvedType {
	return u.types
}

func (u *Union) AddTypes(types ...*ResolvedType) {
	u.types = append(u.types, types...)
}

func (u *Union) SetTypes(types []*ResolvedType) {
	u.types = types
}

func NewUnion(with func(union *Union)) *Union {
	u := &Union{}
	if with != nil {
		with(u)
	}
	return u
}

func (u *Union) typeDef() {}

func (u *Union) Macros() []*MacroInvocation {
	return u.macros
}

func (u *Union) AddMacro(macro *MacroInvocation) {
	u.macros = append(u.macros, macro)
}

func (u *Union) SetMacros(macro []*MacroInvocation) {
	u.macros = macro
}

func (u *Union) Comment() []*CommentLine {
	return u.comment
}

func (u *Union) SetComment(comment []*CommentLine) {
	u.comment = comment
}
