package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"log/slog"
	"path/filepath"
	"strings"
)

type Program struct {
	Pkgs      []*packages.Package
	SSAPkgs   []*ssa.Package
	Program   *wdl.Program
	Callgraph *callgraph.Graph
}

func Parse(dir string) (*Program, error) {
	pkgs, err := packages.Load(
		&packages.Config{
			BuildFlags: []string{"-tags", "macos arm64"},
			Dir:        dir,
			Mode:       packages.NeedDeps | packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedModule | packages.LoadAllSyntax,
		},
		//"...", // this loads all deps and stdlib, but fails and balks with a go mod tidy error
		"./...", // does not load deps and stdlib
	)

	if err != nil {
		return nil, fmt.Errorf("cannot load packages: %w", err)
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

	pg := wdl.NewProgram(nil)

	p := &Program{
		Program:   pg,
		Pkgs:      pkgs,
		SSAPkgs:   ssaPkgs,
		Callgraph: cg,
	}

	return p, p.init()
}

func (p *Program) init() error {
	for _, pkg := range p.Pkgs {
		for _, syntax := range pkg.Syntax {
			for _, decl := range syntax.Decls {
				switch t := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range t.Specs {
						switch spec := spec.(type) {
						case *ast.TypeSpec:
							_, err := p.getTypeDef(p.Program, &wdl.TypeRef{
								Qualifier: wdl.PkgImportQualifier(pkg.PkgPath),
								Name:      wdl.Identifier(spec.Name.Name),
							})

							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (p *Program) getOrInstallFile(pkg *wdl.Package, fname string) *wdl.File {
	for _, file := range pkg.Files() {
		if file.Path() == fname {
			return file
		}
	}

	f := wdl.NewFile(func(file *wdl.File) {
		file.SetName(filepath.Base(fname))
		file.SetPath(filepath.Dir(fname))
	})
	pkg.AddFiles(f)

	return f
}

// getTypeDef inserts (if not yet available) the denoted wdl type and returns it.
func (p *Program) getTypeDef(pg *wdl.Program, ref *wdl.TypeRef) (wdl.TypeDef, error) {
	if ref.Qualifier == "std" {
		switch ref.Name {
		case "Slice":
			orig := p.Program.MustResolveSimple("std", "Slice")
			clone := wdl.NewStruct(func(strct *wdl.Struct) {
				strct.SetPkg(orig.Pkg())
				strct.SetName("Slice")
			})
			tp, err := p.getTypeDef(p.Program, ref.Params[0])
			if err != nil {
				return nil, err
			}
			clone.AddTypeParams(tp.AsResolvedType())
			return clone, nil
		}
	}

	// check short-circuit definition
	if def, ok := pg.TypeDef(ref); ok {
		return def, nil
	}

	// treat as new package and/or type
	var srcPkg *packages.Package
	for _, pkg := range p.Pkgs {
		if pkg.PkgPath == string(ref.Qualifier) {
			srcPkg = pkg
			break
		}
	}

	if srcPkg == nil {
		return nil, fmt.Errorf("src package not found: %s", ref.Qualifier)
	}

	dstPkg, err := p.getOrInstallPackage(ref.Qualifier)
	if err != nil {
		return nil, err
	}

	found := false
	var objPos token.Position
	for ident, object := range srcPkg.TypesInfo.Defs {

		if object == nil {
			continue
		}

		if ident.Name != string(ref.Name) {
			continue
		}

		found = true

		pos := srcPkg.Fset.Position(object.Pos())
		objPos = pos
		file := p.getOrInstallFile(dstPkg, pos.Filename)

		switch obj := object.Type().(type) {
		case *types.Named:
			name := obj.Obj().Name()
			switch obj := obj.Underlying().(type) {
			case *types.Interface:
				// TODO how to distinguish between different use case of a Go interface (type constraint, actually an polymorphic interface etc)
				iface := wdl.NewInterface(func(iface *wdl.Interface) {
					// intentionally add first so that recursion can finish
					dstPkg.AddTypeDefs(iface)
					iface.SetPkg(dstPkg)
					file.AddTypeDefs(iface)

					iface.SetName(wdl.Identifier(name))
					if comment := dstPkg.TypeComments()[iface.Name()]; comment != nil {
						iface.SetComment(comment.Lines())
						iface.SetMacros(comment.Macros())
					}
				})

				if obj.NumMethods() > 0 {

					var methods []*wdl.Func
					//conventional interface methods
					for i := 0; i < obj.NumMethods(); i++ {
						method := obj.Method(i)
						if signature, ok := method.Type().(*types.Signature); ok {
							_ = signature //TODO
						}
						methods = append(methods, wdl.NewFunc(func(fn *wdl.Func) {

							fn.SetName(wdl.Identifier(method.Name()))
							// TODO this is not possible for iface methods and even wrong for global funcs
							if comment := dstPkg.TypeComments()[fn.Name()]; comment != nil {
								fn.SetComment(comment.Lines())
								fn.SetMacros(comment.Macros())
							}
						}))
					}

					return iface, nil
				}

				for i := 0; i < obj.NumEmbeddeds(); i++ {
					switch obj := obj.EmbeddedType(i).(type) {
					case *types.Union:
						// we are a union type definition
						union := wdl.NewUnion(func(union *wdl.Union) {
							// intentionally add first so that recursion can finish
							dstPkg.AddTypeDefs(union)
							union.SetPkg(dstPkg)
							file.AddTypeDefs(union)
							union.SetFile(file)

							union.SetName(wdl.Identifier(name))
							if comment := dstPkg.TypeComments()[union.Name()]; comment != nil {
								union.SetComment(comment.Lines())
								union.SetMacros(comment.Macros())
							}

							for i := 0; i < obj.Len(); i++ {
								ref, err := p.createRef(obj.Term(i).Type())
								if err != nil {
									slog.Error("union: error creating ref", "type", obj.Term(i).Type())
									continue
								}

								tdef, err := p.getTypeDef(p.Program, ref)
								if err != nil {
									slog.Error("union: unsupported term type", slog.String("type", fmt.Sprintf("%T", obj.Term(i).Type())), slog.String("ref", string(union.Name())))
								} else {
									if tdef != nil {
										union.AddTypes(tdef.AsResolvedType())
									}
								}

							}
						})

						return union, nil
					}
				}
			case *types.Struct:
				return wdl.NewStruct(func(strct *wdl.Struct) {
					strct.SetName(wdl.Identifier(name))
					dstPkg.AddTypeDefs(strct)
					strct.SetPkg(dstPkg)
					file.AddTypeDefs(strct)

					if comment := dstPkg.TypeComments()[strct.Name()]; comment != nil {
						strct.SetComment(comment.Lines())
						strct.SetMacros(comment.Macros())
					}

					for fidx := range obj.NumFields() {
						f := obj.Field(fidx)
						strct.AddFields(wdl.NewField(func(field *wdl.Field) {
							field.SetName(wdl.Identifier(f.Name()))
							ref, err := p.createRef(f.Type())
							if err != nil {
								slog.Error("error creating ref for field type", "type", f.Type(), "err", err)
								return
							}
							ftype, err := p.getTypeDef(p.Program, ref)
							if err != nil {
								slog.Error("error getting def for field type", "type", f.Type(), "err", err)
								return
							}

							if ftype == nil {
								slog.Error("oops with nil type for field type", "type", f.Type())
								return
							}

							field.SetTypeDef(ftype.AsResolvedType())
						}))
					}

				}), nil
			case *types.Basic:
				switch obj.Kind() {
				case types.Bool:
					return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
						dType.SetName(wdl.Identifier(name))
						dstPkg.AddTypeDefs(dType)
						dType.SetPkg(dstPkg)
						dType.SetUnderlying(p.Program.MustResolveSimple("std", "bool").TypeDef())
					}), nil
				case types.String:
					return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
						dType.SetName(wdl.Identifier(name))
						dstPkg.AddTypeDefs(dType)
						dType.SetPkg(dstPkg)
						dType.SetUnderlying(p.Program.MustResolveSimple("std", "string").TypeDef())
					}), nil
				case types.Int:
					return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
						dType.SetName(wdl.Identifier(name))
						dstPkg.AddTypeDefs(dType)
						dType.SetPkg(dstPkg)
						dType.SetUnderlying(p.Program.MustResolveSimple("std", "int").TypeDef())
					}), nil
				}
			default:
				slog.Error(fmt.Sprintf("named type not implemented %T@%v", obj, pos))
			}

		case *types.Pointer:
			switch elem := obj.Elem().(type) {
			case *types.Named:
				path := ""
				if elem.Obj().Pkg() != nil {
					path = elem.Obj().Pkg().Path()
				}
				def, err := p.getTypeDef(p.Program, &wdl.TypeRef{
					Qualifier: wdl.PkgImportQualifier(path),
					Name:      wdl.Identifier(elem.Obj().Name()),
				})

				if err != nil {
					return nil, err
				}
				return wdl.NewMutPtr(func(ptr *wdl.MutPtr) {
					ptr.SetTypeDef(def)
				}), nil
			}
			fmt.Printf("%T", obj.Elem())

			//asobj.Elem()
		default:
			slog.Error(fmt.Sprintf("type not implemented %T@%v", obj, objPos))
		}
	}

	slog.Error(fmt.Sprintf("cannot convert def in package %v", ref), "found", found)

	return nil, nil
}

func (p *Program) createRef(typ types.Type) (*wdl.TypeRef, error) {
	switch t := typ.(type) {
	case *types.Named:
		if t.Obj().Pkg() == nil {
			return &wdl.TypeRef{
				Qualifier: "", // this happens for universe types like error
				Name:      wdl.Identifier(t.Obj().Name()),
			}, nil
		}
		return &wdl.TypeRef{
			Qualifier: wdl.PkgImportQualifier(t.Obj().Pkg().Path()),
			Name:      wdl.Identifier(t.Obj().Name()),
		}, nil
	case *types.Pointer:
		ref, err := p.createRef(t.Elem())
		if err != nil {
			return nil, err
		}

		ref.Pointer = true
		return ref, nil
	case *types.Slice:
		ref, err := p.createRef(t.Elem())
		if err != nil {
			return nil, err
		}
		slice := p.Program.MustResolveSimple("std", "Slice").AsTypeRef()
		slice.Params = append(slice.Params, ref)
		return slice, nil
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return p.Program.MustResolveSimple("std", "bool").AsTypeRef(), nil
		case types.String:
			return p.Program.MustResolveSimple("std", "string").AsTypeRef(), nil
		case types.Int:
			return p.Program.MustResolveSimple("std", "int").AsTypeRef(), nil
		case types.Uint:
			return p.Program.MustResolveSimple("std", "uint").AsTypeRef(), nil
		case types.Int32:
			return p.Program.MustResolveSimple("std", "int32").AsTypeRef(), nil
		case types.Int64:
			return p.Program.MustResolveSimple("std", "int64").AsTypeRef(), nil
		case types.Uint32:
			return p.Program.MustResolveSimple("std", "uint32").AsTypeRef(), nil
		case types.Uint64:
			return p.Program.MustResolveSimple("std", "uint64").AsTypeRef(), nil
		case types.Float32:
			return p.Program.MustResolveSimple("std", "float32").AsTypeRef(), nil
		case types.Float64:
			return p.Program.MustResolveSimple("std", "float64").AsTypeRef(), nil
		case types.Byte:
			return p.Program.MustResolveSimple("std", "byte").AsTypeRef(), nil
		}
	}

	return nil, fmt.Errorf("cannot create ref for type %s", typ)
}

// getPackage installs or returns the qualified package.
func (p *Program) getOrInstallPackage(qualifier wdl.PkgImportQualifier) (*wdl.Package, error) {
	res, ok := p.Program.PackageByPath(qualifier)
	if ok {
		return res, nil
	}

	for _, pkg := range p.Pkgs {
		if pkg.PkgPath == string(qualifier) {
			identComments, err := makeIdentComments(pkg)
			if err != nil {
				return nil, err
			}

			res = wdl.NewPackage(func(npkg *wdl.Package) {
				npkg.SetTypeComments(identComments)
				npkg.SetName(wdl.Identifier(pkg.Name))
				npkg.SetQualifier(wdl.PkgImportQualifier(pkg.PkgPath))
			})

			for _, syntax := range pkg.Syntax {
				if syntax.Doc != nil {
					pkgLevelDoc, err := makeComment(pkg, syntax.Doc)
					if err != nil {
						return nil, err
					}

					if res.Comment() == nil {
						res.SetComment(pkgLevelDoc)
					} else {
						res.Comment().AddMacros(pkgLevelDoc.Macros()...)
						res.Comment().AddLines(pkgLevelDoc.Lines()...)
					}
				}
			}

			break
		}
	}

	if res == nil {
		return nil, fmt.Errorf("no such package: %s", qualifier)
	}

	p.Program.AddPackage(res)
	return res, nil
}

func ast2Pos(position token.Position) wdl.Pos {
	return wdl.NewPos(position.Filename, position.Line)
}
