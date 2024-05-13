package wdl

type Param struct {
	name    Identifier
	typeDef *ResolvedType
}

func (p *Param) Name() Identifier {
	return p.name
}

func (p *Param) SetName(name Identifier) {
	p.name = name
}

func (p *Param) TypeDef() *ResolvedType {
	return p.typeDef
}

func (p *Param) SetTypeDef(typeDef *ResolvedType) {
	p.typeDef = typeDef
}

func NewParam(with func(param *Param)) *Param {
	p := &Param{}
	if with != nil {
		with(p)
	}
	return p
}

type Func struct {
	pkg        *Package
	name       Identifier
	args       []*Param
	results    []*Param
	typeParams []*ResolvedType
	receiver   *Param
	visibility Visibility
	body       *BlockStmt
	comment    *Comment
}

func (f *Func) Clone() TypeDef {
	return &Func{
		pkg:        f.pkg,
		name:       f.name,
		comment:    f.comment,
		args:       append([]*Param{}, f.args...),
		results:    append([]*Param{}, f.results...),
		typeParams: append([]*ResolvedType{}, f.typeParams...),
		receiver:   f.receiver,
		visibility: f.visibility,
		body:       f.body,
	}
}

func (f *Func) TypeParams() []*ResolvedType {
	return f.typeParams
}

func (f *Func) SetTypeParams(typeParams []*ResolvedType) {
	f.typeParams = typeParams
}

func (f *Func) AddTypeParams(typeParams ...*ResolvedType) {
	f.typeParams = append(f.typeParams, typeParams...)
}

func (f *Func) Pkg() *Package {
	return f.pkg
}

func (f *Func) SetPkg(pkg *Package) {
	f.pkg = pkg
}

// A Func without a body can be used as a type definition, e.g. for callbacks
func (f *Func) typeDef() {

}

func (f *Func) AsResolvedType() *ResolvedType {
	return NewResolvedType(func(rType *ResolvedType) {
		rType.SetName("")
		rType.SetPkg(f.Pkg())
		rType.SetTypeDef(f)
	})
}

func (f *Func) Body() *BlockStmt {
	return f.body
}

func (f *Func) SetBody(body *BlockStmt) {
	f.body = body
}

func (f *Func) Visibility() Visibility {
	return f.visibility
}

func (f *Func) SetVisibility(visibility Visibility) {
	f.visibility = visibility
}

func (f *Func) Receiver() *Param {
	return f.receiver
}

func (f *Func) SetReceiver(receiver *Param) {
	f.receiver = receiver
}

func (f *Func) Name() Identifier {
	return f.name
}

func (f *Func) SetName(name Identifier) {
	f.name = name
}

func (f *Func) Comment() *Comment {
	return f.comment
}

func (f *Func) SetComment(comment *Comment) {
	f.comment = comment
}

func (f *Func) Args() []*Param {
	return f.args
}

func (f *Func) SetArgs(args []*Param) {
	f.args = args
}

func (f *Func) AddArgs(args ...*Param) {
	f.args = append(f.args, args...)
}

func (f *Func) Results() []*Param {
	return f.results
}

func (f *Func) SetResults(results []*Param) {
	f.results = results
}

func (f *Func) AddResults(results ...*Param) {
	f.results = append(f.results, results...)
}

func NewFunc(with func(fn *Func)) *Func {
	f := &Func{}
	if with != nil {
		with(f)
	}

	return f
}
