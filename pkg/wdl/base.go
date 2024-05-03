package wdl

type Kind int

const (
	TString Kind = iota + 1
	TUInt
	TInt
	TInt32
	TInt64
	TUInt32
	TUInt64
	TFloat32
	TFloat64
	TByte
	TAny
	TBool
)

type BaseType struct {
	pkg  *Package
	name Identifier
	kind Kind
}

func NewBaseType(with func(bt *BaseType)) *BaseType {
	b := &BaseType{}
	if with != nil {
		with(b)
	}

	return b
}

func (b *BaseType) Pkg() *Package {
	return b.pkg
}

func (b *BaseType) SetPkg(pkg *Package) {
	b.pkg = pkg
}

func (b *BaseType) SetName(name Identifier) {
	b.name = name
}

func (b *BaseType) Kind() Kind {
	return b.kind
}

func (b *BaseType) SetKind(base Kind) {
	b.kind = base
}

func (b *BaseType) typeDef() {

}

func (b *BaseType) Name() Identifier {
	return b.name
}

func (b *BaseType) AsResolvedType() *ResolvedType {
	return NewResolvedType(func(rType *ResolvedType) {
		rType.SetName(b.name)
		rType.SetPkg(b.pkg)
		rType.SetTypeDef(b)
	})
}

func (b *BaseType) Macros() []*MacroInvocation {
	return nil
}
