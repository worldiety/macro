package wdl

type Expr interface {
	expression()
}

type BinaryExprOp int

const (
	NEQ BinaryExprOp = iota + 1
	EQL
)

type BinaryExpr struct {
	left     Expr
	operator BinaryExprOp
	right    Expr
}

func (b *BinaryExpr) Left() Expr {
	return b.left
}

func (b *BinaryExpr) SetLeft(left Expr) {
	b.left = left
}

func (b *BinaryExpr) Operator() BinaryExprOp {
	return b.operator
}

func (b *BinaryExpr) SetOperator(operator BinaryExprOp) {
	b.operator = operator
}

func (b *BinaryExpr) Right() Expr {
	return b.right
}

func (b *BinaryExpr) SetRight(right Expr) {
	b.right = right
}

func NewBinaryExpr(with func(bnExpr *BinaryExpr)) *BinaryExpr {
	b := &BinaryExpr{}
	if with != nil {
		with(b)
	}

	return b
}

func (b *BinaryExpr) expression() {}

type StrLit struct {
	value string
}

func (c *StrLit) SetValue(value string) {
	c.value = value
}

func (c *StrLit) Value() string {
	return c.value
}
func NewStrLit(with func(lit *StrLit)) *StrLit {
	c := &StrLit{}
	if with != nil {
		with(c)
	}
	return c
}
func (c *StrLit) expression() {}

type CallExpr struct {
	fun  Expr // function expression, e.g. anon or ident
	args []Expr
}

func (c *CallExpr) Fun() Expr {
	return c.fun
}

func (c *CallExpr) SetFun(fun Expr) {
	c.fun = fun
}

func (c *CallExpr) Args() []Expr {
	return c.args
}

func (c *CallExpr) SetArgs(args []Expr) {
	c.args = args
}

func (c *CallExpr) AddArgs(args ...Expr) {
	c.args = append(c.args, args...)
}

func NewCallExpr(with func(call *CallExpr)) *CallExpr {
	c := &CallExpr{}
	if with != nil {
		with(c)
	}
	return c
}

func (c *CallExpr) expression() {}
