package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
	"reflect"
)

// renderType inspects and emits the actual type.
func (r *Renderer) renderNode(node ast.Node, w *render.BufferedWriter) error {
	switch n := node.(type) {
	case *ast.Struct:
		if err := r.renderStruct(n, w); err != nil {
			return fmt.Errorf("cannot render struct '%s': %w", n.Identifier(), err)
		}
	case *ast.Func:
		return r.renderFunc(n, w)
	case *ast.Block:
		return r.renderBlock(n, w)
	case *ast.Macro:
		if err := r.renderMacro(n, w); err != nil {
			return fmt.Errorf("cannot render macro: %w", err)
		}
	case ast.TypeDecl:
		if err := r.renderTypeDecl(n, w); err != nil {
			return fmt.Errorf("cannot render TypeDecl: %w", err)
		}
	case *ast.CallExpr:
		if err := r.renderCallExpr(n, w); err != nil {
			return fmt.Errorf("cannot render CallExpr: %w", err)
		}
	case *ast.QualIdent:
		if err := r.renderQualIdent(n, w); err != nil {
			return fmt.Errorf("cannot render QualIdent: %w", err)
		}
	case *ast.Ident:
		if err := r.renderIdent(n, w); err != nil {
			return fmt.Errorf("cannot render Ident: %w", err)
		}
	case *ast.SelExpr:
		if err := r.renderSelExpr(n, w); err != nil {
			return fmt.Errorf("cannot render SelExpr: %w", err)
		}
	case *ast.BasicLit:
		if err := r.renderBasicLit(n, w); err != nil {
			return fmt.Errorf("cannot render BasicLit: %w", err)
		}
	case *ast.Assign:
		if err := r.renderAssign(n, w); err != nil {
			return fmt.Errorf("cannot render Assign: %w", err)
		}
	case *ast.Sym:
		if err := r.renderSym(n, w); err != nil {
			return fmt.Errorf("cannot render Sym: %w", err)
		}
	case *ast.IfStmt:
		if err := r.renderIfStmt(n, w); err != nil {
			return fmt.Errorf("cannot render IfStmt: %w", err)
		}
	case *ast.BinaryExpr:
		if err := r.renderBinaryExpr(n, w); err != nil {
			return fmt.Errorf("cannot render BinaryExpr: %w", err)
		}
	case *ast.UnaryExpr:
		if err := r.renderUnaryExpr(n, w); err != nil {
			return fmt.Errorf("cannot render UnaryExpr: %w", err)
		}
	case *ast.ReturnStmt:
		if err := r.renderReturnStmt(n, w); err != nil {
			return fmt.Errorf("cannot render ReturnStmt: %w", err)
		}
	case *ast.CompLit:
		if err := r.renderCompLit(n, w); err != nil {
			return fmt.Errorf("cannot render CompLit: %w", err)
		}

	case *ast.ConstDecl:
		if err := r.renderConst(n, w); err != nil {
			return fmt.Errorf("cannot render ConstDecl: %w", err)
		}
	case *ast.VarDecl:
		if err := r.renderVar(n, w); err != nil {
			return fmt.Errorf("cannot render VarDecl: %w", err)
		}
	case *ast.Interface:
		if err := r.renderInterface(n, w); err != nil {
			return fmt.Errorf("cannot render Interface: %w", err)
		}

	case *ast.Import:
	// handled by Renderer.renderFile

	case *ast.ForStmt:
		if err := r.renderForStmt(n, w); err != nil {
			return fmt.Errorf("cannot render for statement: %w", err)
		}

	case *ast.RangeStmt:
		if err := r.renderRangeStmt(n, w); err != nil {
			return fmt.Errorf("cannot render range statement: %w", err)
		}
	case *ast.Param:
		if err := r.renderParam(n, w); err != nil {
			return fmt.Errorf("cannot render param: %w", err)
		}

	case *ast.DeferStmt:
		if err := r.renderDeferStmt(n, w); err != nil {
			return fmt.Errorf("cannot render defer statement: %w", err)
		}
	case *ast.Tpl:
		if err := r.renderTpl(n, w); err != nil {
			return fmt.Errorf("cannot render template node: %w", err)
		}
	default:
		panic("unsupported type: " + reflect.TypeOf(n).String())
	}

	return nil
}
