package wdl

type RawStmt string

func (b RawStmt) statement() {}

type Statement interface {
	statement()
}

type BlockStmt struct {
	list []Statement
}

func NewBlockStmt(with func(block *BlockStmt)) *BlockStmt {
	b := &BlockStmt{}
	if with != nil {
		with(b)
	}

	return b
}

func (b *BlockStmt) List() []Statement {
	return b.list
}

func (b *BlockStmt) SetList(list []Statement) {
	b.list = list
}

func (b *BlockStmt) Add(list ...Statement) {
	b.list = append(b.list, list...)
}

func (b *BlockStmt) statement() {}

type ReturnStmt struct {
	results []Expr
}

func (r *ReturnStmt) Results() []Expr {
	return r.results
}

func (r *ReturnStmt) SetResults(results []Expr) {
	r.results = results
}

func (r *ReturnStmt) AddResults(results ...Expr) {
	r.results = append(r.results, results...)
}

func NewReturnStmt(with func(*ReturnStmt)) *ReturnStmt {
	r := &ReturnStmt{}
	if with != nil {
		with(r)
	}
	return r
}
func (r *ReturnStmt) statement() {}

type IfStmt struct {
	cond  Expr
	body  *BlockStmt
	initS Statement // may be nil
	elseS Statement // may be nil
}

func (i *IfStmt) Cond() Expr {
	return i.cond
}

func (i *IfStmt) SetCond(cond Expr) {
	i.cond = cond
}

func (i *IfStmt) Body() *BlockStmt {
	return i.body
}

func (i *IfStmt) SetBody(body *BlockStmt) {
	i.body = body
}

func (i *IfStmt) Init() Statement {
	return i.initS
}

func (i *IfStmt) SetInit(initS Statement) {
	i.initS = initS
}

func (i *IfStmt) Else() Statement {
	return i.elseS
}

func (i *IfStmt) SetElse(elseS Statement) {
	i.elseS = elseS
}

func NewIfStmt(with func(ifStmt *IfStmt)) *IfStmt {
	i := &IfStmt{}
	if with != nil {
		with(i)
	}
	return i
}

func (i *IfStmt) statement() {}

type AssignStmt struct {
	lhs []Expr
	rhs []Expr
}

func (a *AssignStmt) statement() {}

func (a *AssignStmt) Lhs() []Expr {
	return a.lhs
}

func (a *AssignStmt) SetLhs(lhs []Expr) {
	a.lhs = lhs
}

func (a *AssignStmt) AddLhs(lhs ...Expr) {
	a.lhs = append(a.lhs, lhs...)
}

func (a *AssignStmt) Rhs() []Expr {
	return a.rhs
}

func (a *AssignStmt) SetRhs(rhs []Expr) {
	a.rhs = rhs
}

func (a *AssignStmt) AddRhs(rhs ...Expr) {
	a.rhs = append(a.rhs, rhs...)
}

func NewAssignStmt(with func(assignStmt *AssignStmt)) *AssignStmt {
	a := &AssignStmt{}
	if with != nil {
		with(a)
	}
	return a
}

type SelectorExpr struct {
	x   Expr
	sel Identifier
}

func (s *SelectorExpr) X() Expr {
	return s.x
}

func (s *SelectorExpr) SetX(x Expr) {
	s.x = x
}

func (s *SelectorExpr) Sel() Identifier {
	return s.sel
}

func (s *SelectorExpr) SetSel(sel Identifier) {
	s.sel = sel
}

func NewSelectorExpr(with func(sel *SelectorExpr)) *SelectorExpr {
	s := &SelectorExpr{}
	if with != nil {
		with(s)
	}
	return s
}

func (s *SelectorExpr) expression() {}
