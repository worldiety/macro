package macro

import (
	"fmt"
	"github.com/worldiety/macro/ast/wdy"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/stdlib"
	"runtime/debug"
	"strings"
)

func (m *GoMacros) TaggedUnion(cfg string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			debug.PrintStack()
		}
	}()
	union := m.typ.(*wdy.Union)
	if !strings.HasPrefix(union.GetRef().Name, "_") {
		panic(fmt.Errorf("invalid union name, must begin with _: %s", union.GetRef()))
	}

	unionName := union.GetRef().Name[1:]
	mod := m.parent.prjMod
	mod.AddPackages(
		ast.NewPkg(union.Ref.Path).
			SetName(union.Ref.PathName()).
			SetPreamble(Preamble).
			AddFiles(
				ast.NewFile(strings.ToLower(unionName + "_gen.go")).
					SetPreamble(Preamble).
					AddTypes(
						ast.NewStruct(unionName).
							SetComment(strings.Join(union.Comment, "\n")).
							AddFields(
								ast.NewField("ordinal", ast.NewSimpleTypeDecl(stdlib.Int)).SetVisibility(ast.Private),
								ast.NewField("value", ast.NewSimpleTypeDecl(stdlib.Any)).SetVisibility(ast.Private),
							),
					),
			),
	)

	return ""
}
