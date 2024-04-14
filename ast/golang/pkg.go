package golang

import (
	"fmt"
	"github.com/worldiety/macro/ast/wdy"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log/slog"
	"regexp"
	"strings"
)

func Load(dir string) ([]*packages.Package, error) {
	pkgs, err := packages.Load(
		&packages.Config{
			Dir:  dir,
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedModule,
		},
		"./...",
	)

	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

var regexMacroCall = regexp.MustCompile(`!{{.+}}`)

func Macros(pkgs []*packages.Package) []wdy.TypeDecl {
	var res []wdy.TypeDecl
	for _, pkg := range pkgs {
		typeDeclrComments := map[string]*ast.CommentGroup{}
		for _, syntax := range pkg.Syntax {
		nextDeclr:
			for _, decl := range syntax.Decls {
				if decl, ok := decl.(*ast.GenDecl); ok {
					if decl.Doc != nil {
						for _, spec := range decl.Specs {
							if spec, ok := spec.(*ast.TypeSpec); ok {
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

		for ident, object := range pkg.TypesInfo.Defs {
			var macros []wdy.Macro
			var commentLines []string
			comment, ok := typeDeclrComments[ident.Name]
			if ok {
				for _, c := range comment.List {
					if macro := regexMacroCall.FindString(c.Text); macro != "" {
						macros = append(macros, wdy.Macro(macro[1:]))
					} else {
						commentLines = append(commentLines, strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(c.Text), "//")))
					}
				}
			}

			if len(commentLines) > 0 && commentLines[len(commentLines)-1] == "" {
				commentLines = commentLines[:len(commentLines)-1]
			}

			if object == nil {
				continue
			}

			switch obj := object.Type().(type) {
			case *types.Named:
				namedRef, _ := intoRef(obj)

				switch obj := obj.Underlying().(type) {
				case *types.Interface:
					for i := 0; i < obj.NumEmbeddeds(); i++ {
						switch obj := obj.EmbeddedType(i).(type) {
						case *types.Union:
							union := &wdy.Union{
								Ref:     namedRef,
								Macros:  macros,
								Comment: commentLines,
							}
							for i := 0; i < obj.Len(); i++ {
								ref, ok := intoRef(obj.Term(i).Type())
								if !ok {
									slog.Error("unsupported term type in union", slog.String("type", fmt.Sprintf("%T", obj.Term(i).Type())), slog.String("ref", union.Ref.String()))
								} else {
									union.Types = append(union.Types, ref)
								}

							}

							res = append(res, union)
						}
					}
				}
			}
		}
	}

	return res
}

func intoRef(typ types.Type) (wdy.TypeReference, bool) {
	switch t := typ.(type) {
	case *types.Basic:
		return wdy.TypeReference{
			Path: "",
			Name: t.Name(),
		}, true
	case *types.Named:
		return wdy.TypeReference{
			Path: t.Obj().Pkg().Path(),
			Name: t.Obj().Name(),
		}, true
	default:
		return wdy.TypeReference{}, false
	}
}
