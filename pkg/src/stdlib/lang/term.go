package lang

import "github.com/worldiety/macro/pkg/src/ast"

// Term writes one or more terminator symbols, e.g.
//
//	Go: \n
//	Java: ; + \n
func Term() *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewSym(ast.SymNewline)),
	)
}
