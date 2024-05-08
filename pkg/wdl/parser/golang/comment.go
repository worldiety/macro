package golang

import (
	"github.com/worldiety/macro/pkg/wdl"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"strings"
)

func makeIdentComments(pkg *packages.Package) (map[wdl.Identifier]*wdl.Comment, error) {
	res := make(map[wdl.Identifier]*wdl.Comment)
	// collect the paired ident and ast comments
	typeDeclrComments := map[string]*ast.CommentGroup{}
	for _, syntax := range pkg.Syntax {
	nextDeclr:
		for _, decl := range syntax.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if decl.Doc != nil {
					typeDeclrComments[decl.Name.Name] = decl.Doc
					continue nextDeclr
				}
			case *ast.GenDecl:
				if decl.Doc != nil {
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
							}

							if spec.Name != nil {
								typeDeclrComments[spec.Name.Name] = decl.Doc
								continue nextDeclr
							}
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

		res[wdl.Identifier(ident)] = wdlComment
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
