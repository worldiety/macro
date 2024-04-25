package fmt

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/stdlib/lang"
)

func Println(args ...ast.Expr) *ast.Macro {
	return lang.CallStatic("fmt.Println", args...)
}
