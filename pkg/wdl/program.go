package wdl

import "fmt"

type Program struct {
	packages []*Package
	std      *Package
	path     string
}

func (p *Program) Path() string {
	return p.path
}

func (p *Program) SetPath(path string) {
	p.path = path
}

func NewProgram(with func(program *Program)) *Program {
	pg := &Program{}
	pg.AddPackage(NewPackage(func(p *Package) {
		pg.std = p
		p.SetName("std")
		p.SetQualifier("std")
		p.SetComment(NewComment(func(comment *Comment) {
			comment.AddLines(NewCommentLine(func(line *CommentLine) {
				line.SetText("std")
			}))
		}))
		p.AddTypeDefs(
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("string")
				bt.SetKind(TString)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("int")
				bt.SetKind(TInt)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("uint")
				bt.SetKind(TUInt)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("int32")
				bt.SetKind(TInt32)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("int64")
				bt.SetKind(TInt64)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("uint32")
				bt.SetKind(TUInt32)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("uint64")
				bt.SetKind(TUInt64)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("float32")
				bt.SetKind(TFloat32)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("float64")
				bt.SetKind(TFloat64)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("byte")
				bt.SetKind(TByte)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("any")
				bt.SetKind(TAny)
			}),
			NewBaseType(func(bt *BaseType) {
				bt.SetPkg(p)
				bt.SetName("bool")
				bt.SetKind(TBool)
			}),
			NewInterface(func(iface *Interface) {
				iface.SetPkg(p)
				iface.SetName("error")
			}),
			NewStruct(func(strct *Struct) {
				strct.SetPkg(p)
				strct.SetName("Slice")
			}),
		)
	}))

	if with != nil {
		with(pg)
	}

	return pg
}

func (p *Program) AddPackage(pkg *Package) {
	p.packages = append(p.packages, pkg)
}

func (p *Program) Packages() []*Package {
	return p.packages
}

func (p *Program) MustResolveSimple(path, ident string) *ResolvedType {
	pkg, ok := p.PackageByPath(PkgImportQualifier(path))
	if !ok {
		panic(fmt.Errorf("could not resolve simple package %s", path))
	}

	for _, def := range pkg.TypeDefs() {
		if def.Name().String() == ident {
			return def.AsResolvedType()
		}
	}

	panic(fmt.Errorf("could not resolve simple type %s.%s", path, ident))
}

func (p *Program) TypeDef(ref *TypeRef) (TypeDef, bool) {
	for _, pkg := range p.packages {
		if pkg.Qualifier() == ref.Qualifier {
			for _, def := range pkg.TypeDefs() {
				if ref.Name == def.Name() {
					return def, true
				}
			}
		}
	}

	return nil, false
}

func (p *Program) PackageByPath(q PkgImportQualifier) (*Package, bool) {
	for _, pkg := range p.packages {
		if pkg.Qualifier() == q {
			return pkg, true
		}
	}

	return nil, false
}
