package golang

import (
	"fmt"
	"github.com/worldiety/macro/ast/wdy"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"log/slog"
	"regexp"
	"strings"
)

type Program struct {
	Pkgs      []*packages.Package
	SSAPkgs   []*ssa.Package
	TypeDecl  []wdy.TypeDecl
	Callgraph *callgraph.Graph
}

func Parse(dir string) (*Program, error) {
	pkgs, err := packages.Load(
		&packages.Config{
			Dir:  dir,
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedModule | packages.LoadAllSyntax,
		},
		"./...",
	)

	if err != nil {
		return nil, err
	}

	// Create and build SSA-form program representation.
	mode := ssa.InstantiateGenerics // instantiate generics by default for soundness
	prog, ssaPkgs := ssautil.AllPackages(pkgs, mode)
	prog.Build()
	goModPath := pkgs[0].Module.Path
	slog.Info("found module path", slog.String("dir", goModPath))

	var cg *callgraph.Graph
	cg = vta.CallGraph(ssautil.AllFunctions(prog), cha.CallGraph(prog))
	cg.DeleteSyntheticNodes()
	err = callgraph.GraphVisitEdges(cg, func(edge *callgraph.Edge) error {
		var callerPkg string
		if edge.Caller.Func.Pkg != nil {
			callerPkg = edge.Caller.Func.Pkg.Pkg.Path()

		}
		callerFuncName := edge.Caller.Func.Name()

		var calleePkg string
		if edge.Callee.Func.Pkg != nil {
			calleePkg = edge.Callee.Func.Pkg.Pkg.Path()

		}
		calleeFuncName := edge.Callee.Func.Name()

		if !strings.HasPrefix(calleePkg, goModPath) && !strings.HasPrefix(callerPkg, goModPath) {
			return nil
		}

		var callerReceiverName string
		if rec := edge.Caller.Func.Signature.Recv(); rec != nil {
			switch t := rec.Type().(type) {
			case *types.Pointer:
				switch t := t.Elem().(type) {
				case *types.Named:
					callerReceiverName = t.Obj().Name()
				}
			case *types.Named:
				callerReceiverName = t.Obj().Name()
			}
		}

		fmt.Printf("%s.%s.%s -> %s.%s\n", callerPkg, callerReceiverName, callerFuncName, calleePkg, calleeFuncName)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("callgraph.GraphVisitEdges: %v", err)
	}

	return &Program{
		Pkgs:      pkgs,
		TypeDecl:  convertTypeDecl(pkgs),
		SSAPkgs:   ssaPkgs,
		Callgraph: cg,
	}, nil
}

var regexMacroCall = regexp.MustCompile(`!{{.+}}`)

func convertTypeDecl(pkgs []*packages.Package) []wdy.TypeDecl {
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
						pos := pkg.Fset.Position(c.Pos())
						macros = append(macros, wdy.Macro{
							Template: macro[1:],
							Origin: wdy.Pos{
								File: pos.Filename,
								Line: pos.Line,
							},
						})
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
					fmt.Println(obj)
					if obj.NumMethods() > 0 {
						var methods []wdy.Func
						//conventional interface methods
						for i := 0; i < obj.NumMethods(); i++ {
							method := obj.Method(i)
							if signature, ok := method.Type().(*types.Signature); ok {
								_ = signature //TODO
							}
							methods = append(methods, wdy.Func{Name: method.Name()})
							// TODO how to access the method doc?
						}
						res = append(res, &wdy.Interface{
							Ref:     namedRef,
							Macros:  macros,
							Comment: commentLines,
							Methods: methods,
						})
					}

					for i := 0; i < obj.NumEmbeddeds(); i++ {
						switch obj := obj.EmbeddedType(i).(type) {
						case *types.Union:
							// we are a union type definition
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
		var path string
		if t.Obj().Pkg() != nil {
			path = t.Obj().Pkg().Path() // e.g. error is no basic type but in universe
		}
		return wdy.TypeReference{
			Path: path,
			Name: t.Obj().Name(),
		}, true
	case *types.Slice:
		tp, ok := intoRef(t.Elem())
		if !ok {
			return wdy.TypeReference{}, false
		}
		return wdy.TypeReference{
			Path:     "",
			Name:     "[]",
			TypeArgs: []wdy.TypeReference{tp},
		}, true
	case *types.Map:
		key, ok := intoRef(t.Key())
		if !ok {
			return wdy.TypeReference{}, false
		}
		val, ok := intoRef(t.Elem())
		if !ok {
			return wdy.TypeReference{}, false
		}
		return wdy.TypeReference{
			Path:     "",
			Name:     "map",
			TypeArgs: []wdy.TypeReference{key, val},
		}, true
	default:
		return wdy.TypeReference{}, false
	}
}
