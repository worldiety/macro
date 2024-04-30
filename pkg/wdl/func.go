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

type Block struct {
	statements []Statement
}

func (b *Block) Statements() []Statement {
	return b.statements
}

func (b *Block) SetStatements(statements []Statement) {
	b.statements = statements
}

func (b *Block) AddStatements(statements ...Statement) {
	b.statements = append(b.statements, statements...)
}

func (b *Block) statement() {}

func NewBlock(with func(blk *Block)) *Block {
	b := &Block{}
	if with != nil {
		with(b)
	}

	return b
}

type RawStmt string

func (b RawStmt) statement() {}

type Statement interface {
	statement()
}

type Func struct {
	pkg        *Package
	name       Identifier
	macros     []*MacroInvocation
	comment    []*CommentLine
	args       []*Param
	results    []*Param
	receiver   *Param
	visibility Visibility
	body       *Block
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

func (f *Func) Body() *Block {
	return f.body
}

func (f *Func) SetBody(body *Block) {
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

func (f *Func) Macros() []*MacroInvocation {
	return f.macros
}

func (f *Func) SetMacros(macros []*MacroInvocation) {
	f.macros = macros
}

func (f *Func) Comment() []*CommentLine {
	return f.comment
}

func (f *Func) SetComment(comment []*CommentLine) {
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
