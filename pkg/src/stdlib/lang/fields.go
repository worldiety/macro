package lang

import "github.com/worldiety/macro/pkg/src/ast"

// Attr returns an expression which refers to a member of the enclosing type of the func.
func Attr(name string) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguageWithContext(ast.LangGo,
			func(m *ast.Macro) []ast.Node {
				var fun *ast.Func
				if ast.ParentAs(m, &fun) {
					recName := fun.RecName()
					if recName == "" {
						if s, ok := fun.Parent().(*ast.Struct); ok {
							recName = s.DefaultRecName
						}
					}

					return ast.Nodes(ast.NewSelExpr(ast.NewIdent(recName), ast.NewIdent(name)))
				}

				return nil
			},
		),
	)
}
