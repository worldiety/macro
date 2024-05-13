package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"go/ast"
	"go/token"
	"log/slog"
	"strconv"
)

func convertStatement(stmt ast.Stmt) (wdl.Statement, error) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		return wdl.NewBlockStmt(func(block *wdl.BlockStmt) {
			for _, s := range stmt.List {
				ws, err := convertStatement(s)
				if err != nil {
					slog.Error("cannot convert statement", "err", err)
					continue
				}

				block.Add(ws)
			}
		}), nil
	case *ast.AssignStmt:
		return wdl.NewAssignStmt(func(assignStmt *wdl.AssignStmt) {
			for _, expr := range stmt.Lhs {
				wlhs, err := convertExpression(expr)
				if err != nil {
					slog.Error("cannot convert assignment lhs node", "err", err)
					continue
				}

				assignStmt.AddLhs(wlhs)
			}

			for _, expr := range stmt.Rhs {
				wrhs, err := convertExpression(expr)
				if err != nil {
					slog.Error("cannot convert assignment rhs node", "err", err)
					continue
				}

				assignStmt.AddRhs(wrhs)
			}
		}), nil
	case *ast.IfStmt:
		return wdl.NewIfStmt(func(ifStmt *wdl.IfStmt) {
			wcond, err := convertExpression(stmt.Cond)
			if err != nil {
				slog.Error("cannot convert if condition", "err", err)
			} else {
				ifStmt.SetCond(wcond)
			}

			if stmt.Init != nil {
				winit, err := convertStatement(stmt.Init)
				if err != nil {
					slog.Error("cannot convert init statement", "err", err)
				} else {
					ifStmt.SetInit(winit)
				}
			}

			if stmt.Else != nil {
				welse, err := convertStatement(stmt.Else)
				if err != nil {
					slog.Error("cannot convert else statement", "err", err)
				} else {
					ifStmt.SetElse(welse)
				}
			}

			wbody, err := convertStatement(stmt.Body)
			if err != nil {
				slog.Error("cannot convert body statement", "err", err)
			} else {
				ifStmt.SetBody(wbody.(*wdl.BlockStmt))
			}
		}), nil
	case *ast.ReturnStmt:
		return wdl.NewReturnStmt(func(wstmt *wdl.ReturnStmt) {
			for _, expr := range stmt.Results {
				we, err := convertExpression(expr)
				if err != nil {
					slog.Error("cannot convert expression", "err", err)
					continue
				}

				wstmt.AddResults(we)
			}
		}), nil
	default:
		return nil, fmt.Errorf("statement not supported: %T", stmt)
	}

}

func convertExpression(expr ast.Expr) (wdl.Expr, error) {
	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		return wdl.NewBinaryExpr(func(bnExpr *wdl.BinaryExpr) {
			wleft, err := convertExpression(expr.X)
			if err != nil {
				slog.Error("cannot convert binary left expression", "err", err)
			}
			bnExpr.SetLeft(wleft)

			switch expr.Op {
			case token.NEQ:
				bnExpr.SetOperator(wdl.NEQ)
			case token.EQL:
				bnExpr.SetOperator(wdl.EQL)
			default:
				slog.Error("unknown binary expr operator", "op", expr.Op.String())
			}

			wright, err := convertExpression(expr.Y)
			if err != nil {
				slog.Error("cannot convert binary right expression", "err", err)
			}
			bnExpr.SetRight(wright)

		}), nil
	case *ast.CallExpr:
		return wdl.NewCallExpr(func(call *wdl.CallExpr) {
			wfun, err := convertExpression(expr.Fun)
			if err != nil {
				slog.Error("cannot convert function expression", "err", err)
			} else {
				call.SetFun(wfun)
			}

			for _, arg := range expr.Args {
				warg, err := convertExpression(arg)
				if err != nil {
					slog.Error("cannot convert argument expression", "err", err)
					continue
				}

				call.AddArgs(warg)
			}

		}), nil
	case *ast.BasicLit:
		switch expr.Kind {
		case token.STRING:
			return wdl.NewStrLit(func(lit *wdl.StrLit) {
				str, err := strconv.Unquote(expr.Value)
				if err != nil {
					slog.Error("cannot unquote literal", "err", err)
				} else {
					lit.SetValue(str)
				}
			}), nil
		default:
			return nil, fmt.Errorf("unknown literal type: %s", expr.Kind.String())
		}
	case *ast.Ident:
		return wdl.Identifier(expr.Name), nil
	case *ast.SelectorExpr:
		return wdl.NewSelectorExpr(func(sel *wdl.SelectorExpr) {
			wx, err := convertExpression(expr.X)
			if err != nil {
				slog.Error("cannot convert x expression on selector", "err", err)
			} else {
				sel.SetX(wx)
			}

			sel.SetSel(wdl.Identifier(expr.Sel.Name))

		}), nil
	default:
		return nil, fmt.Errorf("expr not supported: %T", expr)
	}
}
