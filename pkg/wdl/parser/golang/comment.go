package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"log/slog"
	"strings"
)

func unpackStarExpr(e ast.Expr) string {
	switch exp := e.(type) {
	case *ast.StarExpr:
		return unpackStarExpr(exp.X)
	case *ast.Ident:
		return exp.Name
	case *ast.IndexExpr:
		return unpackStarExpr(exp.X)
	case *ast.IndexListExpr:
		return unpackStarExpr(exp.X)
	}

	panic(e)
}

func makeIdentComments(pkg *packages.Package) (map[wdl.MangeledName]*wdl.Comment, error) {
	res := make(map[wdl.MangeledName]*wdl.Comment)
	// collect the paired ident and ast comments
	typeDeclrComments := map[string]*ast.CommentGroup{}
	for _, syntax := range pkg.Syntax {
	nextDeclr:
		for _, decl := range syntax.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if decl.Doc != nil {
					mangelName := decl.Name.Name
					if decl.Recv != nil && len(decl.Recv.List) > 0 {
						mangelName = unpackStarExpr(decl.Recv.List[0].Type) + "." + mangelName
					}
					typeDeclrComments[mangelName] = decl.Doc
					continue nextDeclr
				}
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.TypeSpec); ok {

						switch def := spec.Type.(type) {
						case *ast.StructType:
							for _, field := range def.Fields.List {
								if field.Doc != nil {
									for _, name := range field.Names {
										typeDeclrComments[spec.Name.Name+"."+name.Name] = field.Doc
									}
								}
							}
						case *ast.InterfaceType:
							for _, method := range def.Methods.List {
								if method.Doc != nil {
									for _, name := range method.Names {
										typeDeclrComments[spec.Name.Name+"."+name.Name] = method.Doc
									}
								}
							}
						default:
							slog.Info("ignored comment genDecl", "type", fmt.Sprintf("%T", spec.Type))
						}

						if spec.Name != nil && decl.Doc != nil {
							typeDeclrComments[spec.Name.Name] = decl.Doc
							continue nextDeclr
						}
					}
				}

			}

		}
	}

	// now transform into comment type
	for ident, comment := range typeDeclrComments {
		wdlComment, err := makeComment(pkg, comment)
		if err != nil {
			return nil, err
		}

		res[wdl.MangeledName(ident)] = wdlComment
	}

	return res, nil
}

func makeComment(pkg *packages.Package, comment *ast.CommentGroup) (*wdl.Comment, error) {
	var macros []*wdl.MacroInvocation
	var commentLines []*wdl.CommentLine

	for _, c := range comment.List {
		if macro := wdl.MacroRegex.FindString(c.Text); macro != "" {
			pos := pkg.Fset.Position(c.Pos())
			m, err := wdl.ParseMacroInvocation(macro, ast2Pos(pos))
			if err != nil {
				return nil, wdl.NewErrorWithPos(ast2Pos(pos), err)
			}
			macros = append(macros, m)
		} else {
			line := wdl.NewCommentLine(func(line *wdl.CommentLine) {
				line.SetPos(ast2Pos(pkg.Fset.Position(c.Pos())))
				line.SetText(strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(c.Text), "//")))
			})
			commentLines = append(commentLines, line)
		}
	}

	if len(commentLines) > 0 && commentLines[len(commentLines)-1].Text() == "" {
		commentLines = commentLines[:len(commentLines)-1]
	}

	return wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(commentLines)
		comment.SetMacros(macros)
	}), nil
}
