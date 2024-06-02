package audit

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
)

const (
	audit = "go.permission.audit"
)

type AddPermissionAnnotation struct {
	prog *wdl.Program
}

func NewAddPermissionAnnotation(prog *wdl.Program) *AddPermissionAnnotation {
	return &AddPermissionAnnotation{prog: prog}
}

func (m *AddPermissionAnnotation) Names() []wdl.MacroName {
	return []wdl.MacroName{
		audit,
	}
}

func (m *AddPermissionAnnotation) Expand(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {
	fn, ok := def.(*wdl.Func)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("macro can only applied to function: %s", audit))
	}

	if fn.Body() == nil || len(fn.Body().List()) == 0 {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("function has empty body"))
	}

	ifStmt, ok := fn.Body().List()[0].(*wdl.IfStmt)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("first statement must be an if-err-guard statement"))
	}

	if ifStmt.Init() == nil {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("if-err-guard statement must use assignment form"))
	}

	// check left-hand-side assignment
	assignStmt, ok := ifStmt.Init().(*wdl.AssignStmt)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("no assignment in if-stmt init: cannot happen"))
	}

	if len(assignStmt.Lhs()) != 1 {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("left-hand-side of assignment in if-err-guard must be exactly one expression"))
	}

	errIdent, ok := assignStmt.Lhs()[0].(wdl.Identifier)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("left-hand-side of assignment in if-err-guard must be an identifier"))
	}

	// check rhs assignment
	if len(assignStmt.Rhs()) != 1 {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("right-hand-side of assignment in if-err-guard must be exactly one expression"))
	}

	auditCall, ok := assignStmt.Rhs()[0].(*wdl.CallExpr)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("right-hand-side of assignment in if-err-guard must be a call expression"))
	}

	auditSelExpr, ok := auditCall.Fun().(*wdl.SelectorExpr)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("right-hand-side of assignment in if-err-guard must be a selector expression"))
	}

	varIdent, ok := auditSelExpr.X().(wdl.Identifier)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("right-hand-side of assignment in if-err-guard must be a selector expression referring to an identifier"))
	}

	// check that varIdent is actually a parameter
	isVarIdent := false
	for _, param := range fn.Args() {
		if param.Name() == varIdent {
			isVarIdent = true
			break
		}
	}

	if !isVarIdent {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid audit call on non-argument identifier"))
	}

	if auditSelExpr.Sel() != "Audit" {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid audit call: expected method '%s' but found '%s'", "Audit", auditSelExpr.Sel()))
	}

	if len(auditCall.Args()) != 1 {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid audit call: expected exact one method parameter"))
	}

	str, ok := auditCall.Args()[0].(*wdl.StrLit)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid audit call: parameter must be a constant string literal"))
	}

	// check if last statement in if-guard returns somehow
	binaryExpr, ok := ifStmt.Cond().(*wdl.BinaryExpr)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid guard call: condition must be a binary expression"))
	}
	if ident, ok := binaryExpr.Left().(wdl.Identifier); !ok || ident != errIdent {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid guard call: left side of condition must be the audit error variable"))
	}

	if ident, ok := binaryExpr.Right().(wdl.Identifier); !ok || ident != "nil" {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid guard call: right side of condition must be nil"))
	}
	if binaryExpr.Operator() != wdl.NEQ {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("invalid guard call: binary operator must be NEQ"))
	}

	if ifStmt.Body() == nil || len(ifStmt.Body().List()) == 0 {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("if-err-guard statement must return but has no statements"))
	}

	ret, ok := ifStmt.Body().List()[len(ifStmt.Body().List())-1].(*wdl.ReturnStmt)
	if !ok {
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("if-err-guard statement must return but last statement in block is not a return"))
	}

	_ = ret // TODO same questions as above

	// finally annotate
	m.prog.AddAnnotations(wdl.NewPermissionAnnotation(func(annotation *wdl.PermissionAnnotation) {
		annotation.SetPermissionID(str.Value())
		annotation.SetTypeDef(def)
	}))

	return nil
}
