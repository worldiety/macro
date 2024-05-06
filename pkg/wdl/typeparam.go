package wdl

type TypeParam struct {
	name Identifier
	pkg  *Package
}

func NewTypeParam(with func(tParm *TypeParam)) *TypeParam {
	t := &TypeParam{}
	if with != nil {
		with(t)
	}
	return t
}

func (t *TypeParam) SetName(name Identifier) {
	t.name = name
}

func (t *TypeParam) Pkg() *Package {
	return t.pkg
}

func (t *TypeParam) SetPkg(pkg *Package) {
	t.pkg = pkg
}

func (t *TypeParam) typeDef() {
}

func (t *TypeParam) Name() Identifier {
	return t.name
}

func (t *TypeParam) AsResolvedType() *ResolvedType {
	return NewResolvedType(func(rType *ResolvedType) {
		rType.SetName(t.name)
		rType.SetPkg(t.pkg)
		rType.SetTypeDef(t)
		rType.SetTypeParam(true)
	})
}

func (t *TypeParam) Macros() []*MacroInvocation {
	return nil
}
