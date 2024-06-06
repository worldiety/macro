package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"log/slog"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
)

type Program struct {
	Pkgs      []*packages.Package
	SSAPkgs   []*ssa.Package
	Program   *wdl.Program
	Callgraph *callgraph.Graph
	//callgraphMap map[wdl.TypeRef]*types.Signature
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
	/*
	   // TODO we cannot use the SSA callgraph representation, because it removed already unreachable code fragments which is common during domain modelling
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

	   		if calleeFuncName == "Audit2" {
	   			constVal := edge.Callee.In[0].Site.(*ssa.Call).Call.Args[0].(*ssa.Const).Value.String()

	   			fmt.Printf("%s.%s.%s -> %s.%s\n", callerPkg, callerReceiverName, callerFuncName, calleePkg, calleeFuncName)
	   			fmt.Println("    ˆ--- ", constVal)
	   		}
	   		_ = callerReceiverName
	   		_ = calleeFuncName
	   		_ = callerFuncName
	   		return nil
	   	})

	   	if err != nil {
	   		return nil, fmt.Errorf("callgraph.GraphVisitEdges: %v", err)
	   	}*/

	pg := wdl.NewProgram(nil)
	if len(pkgs) > 0 {
		pg.SetPath(pkgs[0].Module.Dir)
	}

	p := &Program{
		Program: pg,
		Pkgs:    pkgs,
		SSAPkgs: ssaPkgs,
		//Callgraph: cg,
	}

	return p, p.init()
}

func (p *Program) init() error {
	for _, pkg := range p.Pkgs {
		for _, syntax := range pkg.Syntax {
			for _, decl := range syntax.Decls {
				switch t := decl.(type) {
				case *ast.FuncDecl:
					_, err := p.getTypeDef(p.Program, &wdl.TypeRef{
						Qualifier: wdl.PkgImportQualifier(pkg.PkgPath),
						Name:      wdl.Identifier(t.Name.Name),
					})

					if err != nil {
						return err
					}
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
		if file.AbsolutePath() == fname {
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
func (p *Program) getTypeDef(pg *wdl.Program, ref *wdl.TypeRef) (res wdl.TypeDef, e error) {
	defer func() {
		if res != nil {
			// TODO not sure if we should instantiate generic types within the package or make them virtual?
			if len(ref.Params) > 0 {
				res = res.Clone()
				var tmp []*wdl.ResolvedType
				for _, param := range ref.Params {
					rp, err := p.getTypeDef(pg, param)
					if err != nil {
						e = err
						return
					}
					if rp == nil {
						slog.Error("cannot resolve type parameter", "param", param, "target", ref)
						continue
					}
					tmp = append(tmp, rp.AsResolvedType())
				}
				res.SetTypeParams(tmp)
			}
		}
	}()

	if ref.TypeParam {
		// TODO not sure how to handle that. Seems not to make sense to put that into the package
		return wdl.NewTypeParam(func(tParm *wdl.TypeParam) {
			tParm.SetName(ref.Name)
		}), nil
	}

	/*
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
		}*/

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

		if _, isVar := object.(*types.Var); isVar {
			continue
		}

		if ident.Name != string(ref.Name) {
			continue
		}
		fmt.Println("source type:", srcPkg.PkgPath, ident.Name)
		if object.Name() != ident.Name {
			panic("wtf 1")
		}

		found = true

		pos := srcPkg.Fset.Position(object.Pos())
		objPos = pos
		file := p.getOrInstallFile(dstPkg, pos.Filename)

		objType := object.Type()
		// TODO we have a logical infinite resolver problem between non-pointer ref referring to pointer-types of the same type
		if ptr, ok := objType.(*types.Pointer); ok {
			objType = ptr.Elem()
		}

		switch obj := objType.(type) {
		case *types.Named:
			namedObj := obj
			_ = namedObj
			name := obj.Obj().Name()

			if namedObj.Obj().Name() != ident.Name {
				//panic(fmt.Errorf("wtf 2 %v,%v", object, objType)) TODO fix me: this happens for alias
			}
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
						iface.SetComment(comment)
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
								fn.SetComment(comment)
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
							if strings.HasPrefix(name, "_") {
								union.SetVisibility(visibilityFromName(name[1:]))
							} else {
								union.SetVisibility(visibilityFromName(name))
							}

							union.SetName(wdl.Identifier(name))
							if comment := dstPkg.TypeComments()[union.Name()]; comment != nil {
								union.SetComment(comment)
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

					for tpidx := range namedObj.TypeParams().Len() {
						strct.AddTypeParams(wdl.NewResolvedType(func(rType *wdl.ResolvedType) {
							tp := namedObj.TypeParams().At(tpidx)
							rType.SetTypeParam(true)
							rType.SetName(wdl.Identifier(tp.Obj().Name()))
						}))
					}

					//fmt.Println("!!! CREATED struct", name, dstPkg.Qualifier())
					if ref.Name.String() != name {
						//panic(fmt.Errorf("WTF %s vs %s", ref.Name, name)) TODO fix me this happens for alias
					}
					if comment := dstPkg.TypeComments()[strct.Name()]; comment != nil {
						strct.SetComment(comment)
					}

					for i := range namedObj.NumMethods() {
						method := namedObj.Method(i)
						sig := method.Type().(*types.Signature)
						fn := p.makeFunc(strct.Name().String(), dstPkg, file, srcPkg, method, sig)
						strct.AddMethods(fn)
					}

					for fidx := range obj.NumFields() {
						tag := reflect.StructTag(obj.Tag(fidx))
						f := obj.Field(fidx)
						strct.AddFields(wdl.NewField(func(field *wdl.Field) {
							field.SetName(wdl.Identifier(f.Name()))
							if fieldComment := dstPkg.TypeComments()[strct.Name()+"."+field.Name()]; fieldComment != nil {
								field.SetComment(fieldComment)
							}

							if f.Exported() {
								field.SetVisibility(wdl.Public)
							} else {
								field.SetVisibility(wdl.PackagePrivat)
							}

							value, ok := tag.Lookup("json")
							if ok {
								field.PutTag("json", value)
							}

							constValue, ok := tag.Lookup("value")
							if ok {
								field.PutTag("const", constValue)
							}

							var args []types.Type
							if pnamed, ok := f.Type().(*types.Named); ok {

								if pnamed.TypeArgs().Len() > 0 {
									// this is an instance, e.g. a string
									for i := range pnamed.TypeArgs().Len() {
										args = append(args, pnamed.TypeArgs().At(i))
									}
								} else {
									// this is just a placeholder
									for i := range pnamed.TypeParams().Len() {
										args = append(args, pnamed.TypeArgs().At(i))
									}
								}

							}

							ref, err := p.createRef(f.Type(), args...)
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

							if f.Name() == "Caption" {
								fmt.Println("!!!", ref)
							}

							field.SetTypeDef(ftype.AsResolvedType())
						}))
					}

				}), nil
			case *types.Basic:
				typ, err := p.fromBasicType(dstPkg, obj, name)
				if err != nil {
					return nil, err
				}
				if typ != nil {
					return typ, nil
				}
			case *types.Signature:
				return p.makeFunc("", dstPkg, file, srcPkg, object, obj), nil
			default:
				slog.Error(fmt.Sprintf("named type not implemented %T@%v", obj, pos))
			}

		case *types.Pointer:
		/* TODO this cannot work because we are not using the wdl ref correctly
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
		*/
		//asobj.Elem()
		case *types.Signature:
			// TODO we have them twice, once as a package level but with receiver and once per actual type, we certainly should resolve to a single instance
			return p.makeFunc("", dstPkg, file, srcPkg, object, obj), nil
		default:
			slog.Error(fmt.Sprintf("type not implemented %T@%v", obj, objPos))
		}
	}

	slog.Error(fmt.Sprintf("cannot convert def in package %v", ref), "found", found)

	return nil, nil
}

func (p *Program) makeFunc(receiverTypeName string, dstPkg *wdl.Package, file *wdl.File, srcPkg *packages.Package, nobj types.Object, obj *types.Signature) *wdl.Func {
	return wdl.NewFunc(func(fn *wdl.Func) {
		fn.SetName(wdl.Identifier(nobj.Name()))
		dstPkg.AddTypeDefs(fn)
		fn.SetPkg(dstPkg)
		file.AddTypeDefs(fn)

		if nobj.Exported() {
			fn.SetVisibility(wdl.Public)
		} else {
			fn.SetVisibility(wdl.PackagePrivat)
		}

		var fnDecl *ast.FuncDecl
		for _, syntaxFile := range srcPkg.Syntax {
			for _, decl := range syntaxFile.Decls {
				if astFn, ok := decl.(*ast.FuncDecl); ok {
					if astFn.Name != nil && astFn.Name.Name == nobj.Name() {
						if receiverTypeName != "" && astFn.Recv != nil && len(astFn.Recv.List) > 0 {
							name := ""
							switch def := astFn.Recv.List[0].Type.(type) {
							case *ast.Ident:
								name = def.Name
							case *ast.StarExpr:
								if id, ok := def.X.(*ast.Ident); ok { // todo does not work nested, but thats likely never a domain code anyway
									name = id.Name
								}
							}
							if receiverTypeName != name {
								continue
							}
						}
						fnDecl = astFn
						break
					}
				}
			}
		}

		if fnDecl != nil {
			fn.SetBody(wdl.NewBlockStmt(func(block *wdl.BlockStmt) {
				if fnDecl.Body == nil {
					return // go accepts body-less function signatures, e.g. due to asm implementations
				}
				for _, stmt := range fnDecl.Body.List {
					wstmt, err := convertStatement(stmt)
					if err != nil {
						slog.Error("error converting statements in free func", "stmt", stmt, "err", err, "pkg", dstPkg.Qualifier(), "fn", fn.Name())
						continue
					}
					block.Add(wstmt)
				}
			}))
		}

		for tpidx := range obj.TypeParams().Len() {
			fn.AddTypeParams(wdl.NewResolvedType(func(rType *wdl.ResolvedType) {
				tp := obj.TypeParams().At(tpidx)
				rType.SetTypeParam(true)
				rType.SetName(wdl.Identifier(tp.Obj().Name()))
			}))
		}

		if comment := dstPkg.TypeComments()[fn.Name()]; comment != nil {
			fn.SetComment(comment)
		}

		for i := range obj.Params().Len() {
			param := obj.Params().At(i)
			fn.AddArgs(p.newParam(param))
		}

		for i := range obj.Results().Len() {
			param := obj.Results().At(i)
			fn.AddResults(p.newParam(param))
		}
	})
}

func (p *Program) newParam(varr *types.Var) *wdl.Param {
	return wdl.NewParam(func(param *wdl.Param) {
		param.SetName(wdl.Identifier(varr.Name()))
		ref, err := p.createRef(varr.Type())
		if err != nil {
			slog.Error("error creating ref for param of Var", "type", varr.Type(), "err", err)
			return
		}
		def, err := p.getTypeDef(p.Program, ref)
		if err != nil {
			slog.Error("error getting def for param  of Var", "type", varr.Type(), "err", err)
			return
		}

		if def == nil {
			// panic("unexpected nil definition for param of Var")
			// TODO the resulting ast is wrong and has no type def likely causing weired nil pointers later...
			slog.Error("error getting def for param of Var", "type", varr.Type(), "err", err)
			return
		}
		param.SetTypeDef(def.AsResolvedType())
	})
}

func (p *Program) fromBasicType(dstPkg *wdl.Package, obj *types.Basic, name string) (wdl.TypeDef, error) {
	switch obj.Kind() {
	case types.Bool:
		return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName(wdl.Identifier(name))
			dstPkg.AddTypeDefs(dType)
			dType.SetPkg(dstPkg)
			dType.SetUnderlying(p.Program.MustResolveSimple("std", "bool").TypeDef())

			if comment := dstPkg.TypeComments()[wdl.Identifier(name)]; comment != nil {
				dType.SetComment(comment)
			}
		}), nil
	case types.String:
		return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName(wdl.Identifier(name))
			dstPkg.AddTypeDefs(dType)
			dType.SetPkg(dstPkg)
			dType.SetUnderlying(p.Program.MustResolveSimple("std", "string").TypeDef())

			if comment := dstPkg.TypeComments()[wdl.Identifier(name)]; comment != nil {
				dType.SetComment(comment)
			}
		}), nil

	case types.Float64:
		return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName(wdl.Identifier(name))
			dstPkg.AddTypeDefs(dType)
			dType.SetPkg(dstPkg)
			dType.SetUnderlying(p.Program.MustResolveSimple("std", "float64").TypeDef())

			if comment := dstPkg.TypeComments()[wdl.Identifier(name)]; comment != nil {
				dType.SetComment(comment)
			}
		}), nil
	case types.Float32:
		return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName(wdl.Identifier(name))
			dstPkg.AddTypeDefs(dType)
			dType.SetPkg(dstPkg)
			dType.SetUnderlying(p.Program.MustResolveSimple("std", "float32").TypeDef())

			if comment := dstPkg.TypeComments()[wdl.Identifier(name)]; comment != nil {
				dType.SetComment(comment)
			}
		}), nil
	case types.Int:
		return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName(wdl.Identifier(name))
			dstPkg.AddTypeDefs(dType)
			dType.SetPkg(dstPkg)
			dType.SetUnderlying(p.Program.MustResolveSimple("std", "int").TypeDef())

			if comment := dstPkg.TypeComments()[wdl.Identifier(name)]; comment != nil {
				dType.SetComment(comment)
			}
		}), nil
	case types.Int32:
		return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName(wdl.Identifier(name))
			dstPkg.AddTypeDefs(dType)
			dType.SetPkg(dstPkg)
			dType.SetUnderlying(p.Program.MustResolveSimple("std", "int32").TypeDef())

			if comment := dstPkg.TypeComments()[wdl.Identifier(name)]; comment != nil {
				dType.SetComment(comment)
			}
		}), nil
	case types.Int64:
		return wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName(wdl.Identifier(name))
			dstPkg.AddTypeDefs(dType)
			dType.SetPkg(dstPkg)
			dType.SetUnderlying(p.Program.MustResolveSimple("std", "int64").TypeDef())

			if comment := dstPkg.TypeComments()[wdl.Identifier(name)]; comment != nil {
				dType.SetComment(comment)
			}
		}), nil
	}

	return nil, nil
}

func (p *Program) createRef(typ types.Type, generics ...types.Type) (r *wdl.TypeRef, e error) {
	defer func() {
		if r != nil && len(generics) > 0 {
			// post process generics for all variants, this something between elegant and ugly
			for _, generic := range generics {
				tp, err := p.createRef(generic)
				if err != nil {
					e = err
					return
				}
				r.Params = append(r.Params, tp)
			}
		}
	}()
	switch t := typ.(type) {
	case *types.TypeParam:
		return &wdl.TypeRef{
			Qualifier: "",
			Name:      wdl.Identifier(t.Obj().Name()),
			TypeParam: true,
		}, nil
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
	case *types.Map:
		refKey, err := p.createRef(t.Key())
		if err != nil {
			return nil, err
		}
		refVal, err := p.createRef(t.Elem())
		if err != nil {
			return nil, err
		}
		maep := p.Program.MustResolveSimple("std", "Map").AsTypeRef()
		maep.Params = append(maep.Params, refKey, refVal)
		return maep, nil
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

func visibilityFromName(s string) wdl.Visibility {
	if s == "" {
		return wdl.PackagePrivat
	}

	if strings.HasPrefix(s, "_") {
		return wdl.PackagePrivat
	}

	first, _ := wdl.SplitFirstRune(s)
	if unicode.IsLower(first) {
		return wdl.PackagePrivat
	} else {
		return wdl.Public
	}
}
