package wdl

type Annotation interface {
	isAnnotation()
}

type EntityAnnotation struct {
	alias   string
	pos     Pos
	typeDef TypeDef
}

func (s *EntityAnnotation) Name() string {
	if s.alias != "" {
		return s.alias
	}

	return s.typeDef.Name().String()
}

func (s *EntityAnnotation) Alias() string {
	return s.alias
}

func (s *EntityAnnotation) SetAlias(alias string) {
	s.alias = alias
}

func (s *EntityAnnotation) Pos() Pos {
	return s.pos
}

func (s *EntityAnnotation) SetPos(pos Pos) {
	s.pos = pos
}

func (s *EntityAnnotation) TypeDef() TypeDef {
	return s.typeDef
}

func (s *EntityAnnotation) SetTypeDef(typeDef TypeDef) {
	s.typeDef = typeDef
}

func NewEntityAnnotation(alias string, pos Pos, def TypeDef) *EntityAnnotation {
	return &EntityAnnotation{alias: alias, pos: pos, typeDef: def}
}

func (s *EntityAnnotation) isAnnotation() {}

type BoundedContextAnnotation struct {
	alias string
	pos   Pos
	pkg   *Package
}

func (s *BoundedContextAnnotation) Name() string {
	if s.alias != "" {
		return s.alias
	}

	return s.pkg.Name().String()
}

func (s *BoundedContextAnnotation) Alias() string {
	return s.alias
}

func (s *BoundedContextAnnotation) SetAlias(alias string) {
	s.alias = alias
}

func (s *BoundedContextAnnotation) Pos() Pos {
	return s.pos
}

func (s *BoundedContextAnnotation) SetPos(pos Pos) {
	s.pos = pos
}

func (s *BoundedContextAnnotation) Pkg() *Package {
	return s.pkg
}

func (s *BoundedContextAnnotation) SetPkg(pkg *Package) {
	s.pkg = pkg
}

func NewBoundedContextAnnotation(alias string, pkg *Package, pos Pos) *BoundedContextAnnotation {
	return &BoundedContextAnnotation{alias: alias, pos: pos, pkg: pkg}
}

func (s *BoundedContextAnnotation) isAnnotation() {}

type UseCaseAnnotation struct {
	alias string
	pos   Pos
	fn    *Func
}

func (s *UseCaseAnnotation) Name() string {
	if s.alias != "" {
		return s.alias
	}

	return s.fn.Name().String()
}

func (s *UseCaseAnnotation) Alias() string {
	return s.alias
}

func (s *UseCaseAnnotation) SetAlias(alias string) {
	s.alias = alias
}

func (s *UseCaseAnnotation) Pos() Pos {
	return s.pos
}

func (s *UseCaseAnnotation) SetPos(pos Pos) {
	s.pos = pos
}

func (s *UseCaseAnnotation) Fn() *Func {
	return s.fn
}

func (s *UseCaseAnnotation) SetFn(fn *Func) {
	s.fn = fn
}

func NewUseCaseAnnotation(alias string, fn *Func, pos Pos) *UseCaseAnnotation {
	return &UseCaseAnnotation{alias: alias, fn: fn, pos: pos}
}

func (s *UseCaseAnnotation) isAnnotation() {}

type RepositoryAnnotation struct {
	alias   string
	pos     Pos
	typeDef TypeDef
}

func NewRepositoryAnnotation(alias string, pos Pos, def TypeDef) *RepositoryAnnotation {
	return &RepositoryAnnotation{alias: alias, pos: pos, typeDef: def}
}

func (s *RepositoryAnnotation) isAnnotation() {}

type AggregateRootAnnotation struct {
	alias   string
	pos     Pos
	typeDef TypeDef
}

func (s *AggregateRootAnnotation) Name() string {
	if s.alias != "" {
		return s.alias
	}

	return s.typeDef.Name().String()
}

func (s *AggregateRootAnnotation) Alias() string {
	return s.alias
}

func (s *AggregateRootAnnotation) SetAlias(alias string) {
	s.alias = alias
}

func (s *AggregateRootAnnotation) Pos() Pos {
	return s.pos
}

func (s *AggregateRootAnnotation) SetPos(pos Pos) {
	s.pos = pos
}

func (s *AggregateRootAnnotation) TypeDef() TypeDef {
	return s.typeDef
}

func (s *AggregateRootAnnotation) SetTypeDef(typeDef TypeDef) {
	s.typeDef = typeDef
}

func NewAggregateRootAnnotation(alias string, pos Pos, def TypeDef) *AggregateRootAnnotation {
	return &AggregateRootAnnotation{alias: alias, pos: pos, typeDef: def}
}

func (s *AggregateRootAnnotation) isAnnotation() {}

type ValueAnnotation struct {
	alias   string
	pos     Pos
	typeDef TypeDef
}

func (s *ValueAnnotation) Name() string {
	if s.alias != "" {
		return s.alias
	}

	return s.typeDef.Name().String()
}

func (s *ValueAnnotation) Alias() string {
	return s.alias
}

func (s *ValueAnnotation) SetAlias(alias string) {
	s.alias = alias
}

func (s *ValueAnnotation) Pos() Pos {
	return s.pos
}

func (s *ValueAnnotation) SetPos(pos Pos) {
	s.pos = pos
}

func (s *ValueAnnotation) TypeDef() TypeDef {
	return s.typeDef
}

func (s *ValueAnnotation) SetTypeDef(typeDef TypeDef) {
	s.typeDef = typeDef
}

func NewValueAnnotation(alias string, pos Pos, def TypeDef) *ValueAnnotation {
	return &ValueAnnotation{alias: alias, pos: pos, typeDef: def}
}

func (s *ValueAnnotation) isAnnotation() {}

type DomainEventAnnotation struct {
	alias   string
	pos     Pos
	typeDef TypeDef
}

func (s *DomainEventAnnotation) Name() string {
	if s.alias != "" {
		return s.alias
	}

	return s.typeDef.Name().String()
}

func (s *DomainEventAnnotation) Alias() string {
	return s.alias
}

func (s *DomainEventAnnotation) SetAlias(alias string) {
	s.alias = alias
}

func (s *DomainEventAnnotation) Pos() Pos {
	return s.pos
}

func (s *DomainEventAnnotation) SetPos(pos Pos) {
	s.pos = pos
}

func (s *DomainEventAnnotation) TypeDef() TypeDef {
	return s.typeDef
}

func (s *DomainEventAnnotation) SetTypeDef(typeDef TypeDef) {
	s.typeDef = typeDef
}

func NewDomainEventAnnotation(alias string, pos Pos, def TypeDef) *DomainEventAnnotation {
	return &DomainEventAnnotation{alias: alias, pos: pos, typeDef: def}
}

func (s *DomainEventAnnotation) isAnnotation() {}

type DomainServiceAnnotation struct {
	alias   string
	pos     Pos
	typeDef TypeDef
}

func (s *DomainServiceAnnotation) Alias() string {
	return s.alias
}

func (s *DomainServiceAnnotation) SetAlias(alias string) {
	s.alias = alias
}

func (s *DomainServiceAnnotation) Pos() Pos {
	return s.pos
}

func (s *DomainServiceAnnotation) SetPos(pos Pos) {
	s.pos = pos
}

func (s *DomainServiceAnnotation) TypeDef() TypeDef {
	return s.typeDef
}

func (s *DomainServiceAnnotation) SetTypeDef(typeDef TypeDef) {
	s.typeDef = typeDef
}

func NewDomainServiceAnnotation(alias string, pos Pos, def TypeDef) *DomainServiceAnnotation {
	return &DomainServiceAnnotation{alias: alias, pos: pos, typeDef: def}
}

func (s *DomainServiceAnnotation) isAnnotation() {}
func (s *DomainServiceAnnotation) Name() string {
	if s.alias != "" {
		return s.alias
	}

	return s.typeDef.Name().String()
}
