package wdl

import "slices"

type Visibility int

const (
	PackagePrivat Visibility = iota
	Public
)

type Field struct {
	name       Identifier
	typeDef    *ResolvedType
	visibility Visibility
	tags       map[string]string
}

func (f *Field) PutTag(key, value string) {
	f.tags[key] = value
}

func (f *Field) Tags() map[string]string {
	return f.tags
}

func (f *Field) SetTags(tags map[string]string) {
	f.tags = tags
}

func (f *Field) Visibility() Visibility {
	return f.visibility
}

func (f *Field) SetVisibility(visibility Visibility) {
	f.visibility = visibility
}

func NewField(with func(field *Field)) *Field {
	f := &Field{tags: make(map[string]string)}
	if with != nil {
		with(f)
	}

	return f
}

func (f *Field) Name() Identifier {
	return f.name
}

func (f *Field) SetName(name Identifier) {
	f.name = name
}

func (f *Field) TypeDef() *ResolvedType {
	return f.typeDef
}

func (f *Field) SetTypeDef(typeDef *ResolvedType) {
	f.typeDef = typeDef
}

type Struct struct {
	pkg        *Package
	macros     []*MacroInvocation
	comment    []*CommentLine
	types      []*ResolvedType // composition?
	name       Identifier
	fields     []*Field
	methods    []*Func
	visibility Visibility
	typeParams []*ResolvedType
}

func (s *Struct) TypeParams() []*ResolvedType {
	return s.typeParams
}

func (s *Struct) SetTypeParams(typeParams []*ResolvedType) {
	s.typeParams = typeParams
}

func (s *Struct) AddTypeParams(typeParams ...*ResolvedType) {
	s.typeParams = append(s.typeParams, typeParams...)
}

func (s *Struct) Visibility() Visibility {
	return s.visibility
}

func (s *Struct) SetVisibility(visibility Visibility) {
	s.visibility = visibility
}

func (s *Struct) Methods() []*Func {
	return s.methods
}

func (s *Struct) SetMethods(methods []*Func) {
	s.methods = methods
}

func (s *Struct) AddMethods(methods ...*Func) {
	s.methods = append(s.methods, methods...)
}

func NewStruct(with func(strct *Struct)) *Struct {
	s := &Struct{}
	if with != nil {
		with(s)
	}

	return s
}

func (s *Struct) Pkg() *Package {
	return s.pkg
}

func (s *Struct) SetPkg(pkg *Package) {
	s.pkg = pkg
}

func (s *Struct) Macros() []*MacroInvocation {
	return s.macros
}

func (s *Struct) SetMacros(macros []*MacroInvocation) {
	s.macros = macros
}

func (s *Struct) Comment() []*CommentLine {
	return s.comment
}

func (s *Struct) SetComment(comment []*CommentLine) {
	s.comment = comment
}

func (s *Struct) Types() []*ResolvedType {
	return s.types
}

func (s *Struct) SetTypes(types []*ResolvedType) {
	s.types = types
}

func (s *Struct) Name() Identifier {
	return s.name
}

func (s *Struct) SetName(name Identifier) {
	s.name = name
}

func (s *Struct) Fields() []*Field {
	return s.fields
}

func (s *Struct) SetFields(fields []*Field) {
	s.fields = fields
}

func (s *Struct) AddFields(fields ...*Field) {
	s.fields = append(s.fields, fields...)
}

func (s *Struct) typeDef() {}

func (s *Struct) AsResolvedType() *ResolvedType {
	return NewResolvedType(func(rType *ResolvedType) {
		rType.SetName(s.name)
		rType.SetPkg(s.pkg)
		rType.SetTypeDef(s)
		rType.SetParams(slices.Clone(s.TypeParams()))
	})
}
