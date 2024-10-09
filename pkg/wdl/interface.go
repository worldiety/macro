package wdl

type Interface struct {
	pkg        *Package
	name       Identifier
	types      []*ResolvedType // todo is this also union types?
	typeParams []*ResolvedType
	comment    *Comment
	visibility Visibility
	methods    []*Func
}

func (u *Interface) Clone() TypeDef {
	return &Interface{
		pkg:        u.pkg,
		name:       u.name,
		comment:    u.comment,
		types:      append([]*ResolvedType{}, u.types...),
		typeParams: append([]*ResolvedType{}, u.typeParams...),
	}
}

func (u *Interface) SetTypeParams(typeParams []*ResolvedType) {
	u.typeParams = typeParams
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

func (u *Interface) Comment() *Comment {
	return u.comment
}

func (u *Interface) SetComment(comment *Comment) {
	u.comment = comment
}

func (u *Interface) Visibility() Visibility {
	return u.visibility
}

func (u *Interface) SetVisibility(visibility Visibility) {
	u.visibility = visibility
}

func (u *Interface) Methods() []*Func {
	return u.methods
}

func (u *Interface) SetMethods(methods []*Func) {
	u.methods = methods
}

func (u *Interface) AddMethod(method *Func) {
	u.methods = append(u.methods, method)
}
