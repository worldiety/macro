package lang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/stdlib"
)

// CallDefine emits a variable (re)declaration with an assignment.
func CallDefine(lhs, rhs ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewAssign(ast.Exprs(lhs), ast.AssignDefine, ast.Exprs(rhs))),
	)
}

// TryDefine emits a variable (re)declaration with an assignment and an error check with early return.
// It evaluates the current context to decide how to return and how to re-throw error.
func TryDefine(lhs, rhs ast.Expr, errMsg string) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguageWithContext(ast.LangGo,
			func(m *ast.Macro) []ast.Node {
				myFunc := assertFunc(m)
				if len(myFunc.FunResults) == 0 {
					panic("func " + myFunc.FunName + " must define at least an error return value")
				}

				lastSTD, ok := myFunc.FunResults[len(myFunc.FunResults)-1].ParamTypeDecl.(*ast.SimpleTypeDecl)
				if !ok || lastSTD.SimpleName != stdlib.Error {
					panic("func " + myFunc.FunName + " last result must be an error return value but is " + fmt.Sprint(lastSTD))
				}

				var results []ast.Expr
				for i := 0; i < len(myFunc.FunResults)-1; i++ {
					decl := myFunc.FunResults[i].TypeDecl()
					switch t := decl.(type) {
					case *ast.SimpleTypeDecl:
						if stdlib.IsNumber(string(t.SimpleName)) {
							results = append(results, ast.NewIntLit(0))
						} else if t.SimpleName == stdlib.String {
							results = append(results, ast.NewStrLit(""))
						} else {
							if isInterface(myFunc, t.SimpleName) {
								results = append(results, ast.NewIdent("nil"))
							} else {
								// TODO this is not always correct and cannot always be resolved due to external dependencies
								results = append(results, ast.NewCompLit(t.Clone()))
							}
						}
					case *ast.SliceTypeDecl:
						results = append(results, ast.NewIdent("nil"))
					case *ast.TypeDeclPtr:
						results = append(results, ast.NewIdent("nil"))
					}

				}

				results = append(results, CallStatic("fmt.Errorf", ast.NewStrLit(errMsg+": %w"), ast.NewIdent("err")))

				if lhs == nil {
					return ast.Nodes(
						ast.NewIfStmt(ast.NewBinaryExpr(ast.NewIdent("err"), ast.OpNotEqual, ast.NewIdent("nil")), ast.NewBlock(
							ast.NewReturnStmt(results...),
						)).SetInit(ast.NewAssign(ast.Exprs(lhs, ast.NewIdent("err")), ast.AssignDefine, ast.Exprs(rhs))),

						Term(),
					)
				} else {
					return ast.Nodes(
						ast.NewAssign(ast.Exprs(lhs, ast.NewIdent("err")), ast.AssignDefine, ast.Exprs(rhs)),
						Term(),
						ast.NewIfStmt(ast.NewBinaryExpr(ast.NewIdent("err"), ast.OpNotEqual, ast.NewIdent("nil")), ast.NewBlock(
							ast.NewReturnStmt(results...),
						)),
						Term(),
					)
				}
			},
		),
	)
}

// not yet correct but can identify package local interfaces
func isInterface(scope ast.Node, name ast.Name) bool {
	pkg := assertPkg(scope)

	for _, i := range pkg.Interfaces() {
		if i.TypeName == string(name) {
			return true
		}
	}

	return false
}

// there is always an outer func definition
func assertFunc(n ast.Node) *ast.Func {
	f := &ast.Func{}
	if ok := ast.ParentAs(n, &f); ok {
		return f
	}

	panic("invalid context: must be a func child")
}

// there is always an outer pkg definition
func assertPkg(n ast.Node) *ast.Pkg {
	f := &ast.Pkg{}
	if ok := ast.ParentAs(n, &f); ok {
		return f
	}

	panic("invalid context: must be a pkg child")
}
