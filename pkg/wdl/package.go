package wdl

import "maps"

// A PkgName is the identifier of a package. Examples:
//   - Go: http, ast2, math (this has nothing to do with last segment of the import path)
//   - Java: util, collections (always the last segment)
//   - TypeScript: there is no package system. Identifiers are imported directly from file imports (just the relative path without file extension).
//   - Rust: std, io (looks like there are multiple ways to define them, some conventional by path, some explicit)
type PkgName Identifier

// A PkgImportQualifier denotes the unique and distinct qualifier for a package.
// Examples:
//   - Go: encoding/json (a path)
//   - Java: java.util
//   - TypeScript: @/shared/protocol/gen/event (a path)
//   - Rust: std::io
type PkgImportQualifier string

type Package struct {
	name         Identifier
	qualifier    PkgImportQualifier
	typeDefs     []TypeDef
	typeComments map[Identifier]*Comment
	comment      *Comment
	files        []*File
}

func (p *Package) Pkg() *Package {
	return p
}

func (p *Package) typeDef() {
}

func (p *Package) AsResolvedType() *ResolvedType {
	panic("package is a fake type")
}

func (p *Package) Macros() []*MacroInvocation {
	if p.comment != nil {
		return p.comment.Macros()
	}

	return nil
}

func (p *Package) Clone() TypeDef {
	return &Package{
		name:         p.name,
		qualifier:    p.qualifier,
		typeDefs:     append([]TypeDef(nil), p.typeDefs...),
		typeComments: maps.Clone(p.typeComments),
		comment:      p.comment,
		files:        append([]*File{}, p.files...),
	}
}

func (p *Package) SetTypeParams(typeParams []*ResolvedType) {
	//TODO implement me
	panic("implement me")
}

func (p *Package) Files() []*File {
	return p.files
}

func (p *Package) SetFiles(files []*File) {
	for _, file := range files {
		file.SetPkg(p)
	}
	p.files = files
}

func (p *Package) AddFiles(files ...*File) {
	for _, file := range files {
		file.SetPkg(p)
	}
	p.files = append(p.files, files...)
}

func (p *Package) TypeComments() map[Identifier]*Comment {
	return p.typeComments
}

func (p *Package) SetTypeComments(typeComments map[Identifier]*Comment) {
	p.typeComments = typeComments
}

func (p *Package) Comment() *Comment {
	return p.comment
}

func (p *Package) SetComment(comment *Comment) {
	p.comment = comment
}

func NewPackage(with func(pkg *Package)) *Package {
	pkg := &Package{}
	if with != nil {
		with(pkg)
	}

	return pkg
}

func (p *Package) Name() Identifier {
	return p.name
}

func (p *Package) SetName(name Identifier) {
	p.name = name
}

func (p *Package) Qualifier() PkgImportQualifier {
	return p.qualifier
}

func (p *Package) SetQualifier(qualifier PkgImportQualifier) {
	p.qualifier = qualifier
}

func (p *Package) TypeDefs() []TypeDef {
	return p.typeDefs
}

func (p *Package) SetTypeDefs(typeDefs []TypeDef) {
	p.typeDefs = typeDefs
}

func (p *Package) AddTypeDefs(typeDefs ...TypeDef) {
	p.typeDefs = append(p.typeDefs, typeDefs...)
}
