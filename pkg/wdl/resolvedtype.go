package wdl

// ResolvedType tries to simplify things.
type ResolvedType struct {
	pkg     *Package
	name    Identifier
	typeDef TypeDef
	params  []*ResolvedType
	pointer bool
}

func (r *ResolvedType) Pkg() *Package {
	return r.pkg
}

func (r *ResolvedType) SetPkg(pkg *Package) {
	r.pkg = pkg
}

func (r *ResolvedType) Name() Identifier {
	return r.name
}

func (r *ResolvedType) SetName(name Identifier) {
	r.name = name
}

func (r *ResolvedType) TypeDef() TypeDef {
	return r.typeDef
}

func (r *ResolvedType) SetTypeDef(typeDef TypeDef) {
	r.typeDef = typeDef
}

func (r *ResolvedType) Params() []*ResolvedType {
	return r.params
}

func (r *ResolvedType) SetParams(params []*ResolvedType) {
	r.params = params
}

func (r *ResolvedType) AddParams(params ...*ResolvedType) {
	r.params = append(r.params, params...)
}

func (r *ResolvedType) Pointer() bool {
	return r.pointer
}

func (r *ResolvedType) SetPointer(pointer bool) {
	r.pointer = pointer
}

func NewResolvedType(with func(rType *ResolvedType)) *ResolvedType {
	t := &ResolvedType{}
	if with != nil {
		with(t)
	}

	return t
}
