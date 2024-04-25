package strings

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/stdlib/lang"
)

func NewStrBuilder(ident string, writeStrings ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguageWithContext(ast.LangGo,
			func(m *ast.Macro) []ast.Node {
				decl := ast.NewAssign(ast.Exprs(ast.NewIdent(ident)), ast.AssignDefine, ast.Exprs(ast.NewUnaryExpr(lang.CreateLiteral("strings.Builder"), ast.OpAnd)))
				var nodes []ast.Node
				nodes = append(nodes, decl, ast.NewSym(ast.SymNewline))
				for _, writeString := range writeStrings {
					nodes = append(nodes, lang.CallIdent("sb", "WriteString", writeString), ast.NewSym(ast.SymNewline))
				}

				return nodes
			},
		),
	)
}
