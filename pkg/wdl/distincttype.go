package wdl

import "slices"

type DistinctType struct {
	pkg        *Package
	macros     []*MacroInvocation
	comment    []*CommentLine
	types      []*ResolvedType
	name       Identifier
	visibility Visibility
	typeParams []*ResolvedType
	underlying TypeDef
	methods    []*Func
}

func (d *DistinctType) Clone() TypeDef {
	return &DistinctType{
		pkg:        d.pkg,
		macros:     append([]*MacroInvocation{}, d.macros...),
		comment:    append([]*CommentLine{}, d.comment...),
		types:      append([]*ResolvedType{}, d.types...),
		name:       d.name,
		visibility: d.visibility,
		typeParams: append([]*ResolvedType{}, d.typeParams...),
		underlying: d.underlying,
		methods:    append([]*Func{}, d.methods...),
	}
}

func (d *DistinctType) Methods() []*Func {
	return d.methods
}

func (d *DistinctType) SetMethods(methods []*Func) {
	d.methods = methods
}

func (d *DistinctType) typeDef() {
}

func (d *DistinctType) AsResolvedType() *ResolvedType {
	return &ResolvedType{
		pkg:     d.pkg,
		name:    d.name,
		typeDef: d,
		params:  slices.Clone(d.typeParams),
	}
}

func (d *DistinctType) Pkg() *Package {
	return d.pkg
}

func (d *DistinctType) SetPkg(pkg *Package) {
	d.pkg = pkg
}

func (d *DistinctType) Macros() []*MacroInvocation {
	return d.macros
}

func (d *DistinctType) SetMacros(macros []*MacroInvocation) {
	d.macros = macros
}

func (d *DistinctType) Comment() []*CommentLine {
	return d.comment
}

func (d *DistinctType) SetComment(comment []*CommentLine) {
	d.comment = comment
}

func (d *DistinctType) Types() []*ResolvedType {
	return d.types
}

func (d *DistinctType) SetTypes(types []*ResolvedType) {
	d.types = types
}

func (d *DistinctType) Name() Identifier {
	return d.name
}

func (d *DistinctType) SetName(name Identifier) {
	d.name = name
}

func (d *DistinctType) Visibility() Visibility {
	return d.visibility
}

func (d *DistinctType) SetVisibility(visibility Visibility) {
	d.visibility = visibility
}

func (d *DistinctType) TypeParams() []*ResolvedType {
	return d.typeParams
}

func (d *DistinctType) SetTypeParams(typeParams []*ResolvedType) {
	d.typeParams = typeParams
}

func (d *DistinctType) Underlying() TypeDef {
	return d.underlying
}

func (d *DistinctType) SetUnderlying(underlying TypeDef) {
	d.underlying = underlying
}

func NewDistinctType(with func(dType *DistinctType)) *DistinctType {
	d := &DistinctType{}
	if with != nil {
		with(d)
	}
	return d
}
