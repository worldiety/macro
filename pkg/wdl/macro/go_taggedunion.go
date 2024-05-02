package macro

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render/golang"
	"strings"
)

type jsonRepr string

const (
	internallyTagged jsonRepr = "intern"
)

type goTaggedUnionParams struct {
	JSONRepresentation jsonRepr `json:"json"` // currently only "intern"
	TagName            string   `json:"tag"`
}

func (e *Engine) goTaggedUnion(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {
	if _, ok := def.(*wdl.Interface); ok {
		// TODO this case happens, cannot decide properly if that is correct or wrong
		return nil
	}
	union, ok := def.(*wdl.Union)
	if !ok {
		return fmt.Errorf("expected union definition")
	}

	var opts goTaggedUnionParams
	if err := macroInvoc.UnmarshalParams(&opts); err != nil {
		return fmt.Errorf("invalid macro params: %w", err)
	}

	if opts.TagName == "" {
		opts.TagName = "type"
	}

	if opts.JSONRepresentation == "" {
		opts.JSONRepresentation = internallyTagged
	}

	uStruct := wdl.NewStruct(func(strct *wdl.Struct) {
		strct.SetPkg(union.Pkg())
		strct.SetVisibility(wdl.Public)
		strct.SetName(wdl.Identifier(stripName(union.Name().String())))
		strct.SetComment(union.Comment())
		strct.AddFields(
			wdl.NewField(func(field *wdl.Field) {
				field.SetName("ordinal")
				field.SetTypeDef(e.prog.MustResolveSimple("std", "int"))
			}),
			wdl.NewField(func(field *wdl.Field) {
				field.SetName("value")
				field.SetTypeDef(e.prog.MustResolveSimple("std", "any"))
			}),
		)

		strct.AddMethods(
			wdl.NewFunc(func(fn *wdl.Func) {
				fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
					r.SetName("e")
					r.SetTypeDef(strct.AsResolvedType())
				}))
				fn.SetVisibility(wdl.Public)
				fn.SetName("Unwrap")
				fn.AddResults(
					wdl.NewParam(func(param *wdl.Param) {
						param.SetTypeDef(e.prog.MustResolveSimple("std", "any"))
					}),
				)
				fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {

					blk.AddStatements(wdl.RawStmt(fmt.Sprintf("return e.value")))
				}))
			}),
			wdl.NewFunc(func(fn *wdl.Func) {
				fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
					r.SetName("e")
					r.SetTypeDef(strct.AsResolvedType())
				}))
				fn.SetVisibility(wdl.Public)
				fn.SetName("Ordinal")
				fn.AddResults(
					wdl.NewParam(func(param *wdl.Param) {
						param.SetTypeDef(e.prog.MustResolveSimple("std", "int"))
					}),
				)
				fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {

					blk.AddStatements(wdl.RawStmt(fmt.Sprintf("return e.ordinal")))
				}))
			}),
			wdl.NewFunc(func(fn *wdl.Func) {
				fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
					r.SetName("e")
					r.SetTypeDef(strct.AsResolvedType())
				}))
				fn.SetVisibility(wdl.Public)
				fn.SetName("Valid")
				fn.AddResults(
					wdl.NewParam(func(param *wdl.Param) {
						param.SetTypeDef(e.prog.MustResolveSimple("std", "bool"))
					}),
				)
				fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {

					blk.AddStatements(wdl.RawStmt(fmt.Sprintf("return e.ordinal>0")))
				}))
			}),
			wdl.NewFunc(func(fn *wdl.Func) {
				fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
					r.SetName("e")
					r.SetTypeDef(strct.AsResolvedType())
				}))
				fn.SetVisibility(wdl.Public)
				fn.SetName("Switch")
				fn.SetComment(wdl.NewSimpleComment("Switch provides an exhaustive and type safe closure callback mechanic. Nil callbacks are allowed. Unmatched branches are delegated into a default case.").Lines())
				for _, resolvedType := range union.Types() {
					fn.AddArgs(
						wdl.NewParam(func(param *wdl.Param) {
							param.SetName("on" + identFrom(resolvedType))
							param.SetTypeDef(wdl.NewFunc(func(fn *wdl.Func) {
								fn.SetPkg(union.Pkg())
								fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
									param.SetTypeDef(resolvedType)
								}))
							}).AsResolvedType())
						}),
					)
				}

				fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
					param.SetName("_onDefault") // avoid accidental name collisions using prefix
					param.SetTypeDef(wdl.NewFunc(func(fn *wdl.Func) {
						fn.SetPkg(union.Pkg())
						fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
							param.SetTypeDef(e.prog.MustResolveSimple("std", "any"))
						}))
					}).AsResolvedType())
				}))

				fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {
					tmp := "switch e.ordinal {\n"
					for idx, resolvedType := range union.Types() {
						ord := idx + 1
						rtmp := golang.NewRFile(golang.NewRenderer(golang.Options{}), union.Pkg().Qualifier())
						gtype := rtmp.GoType(resolvedType)

						tmp += fmt.Sprintf("case %d:\n", ord)
						tmp += fmt.Sprintf("if on%s !=nil {\n", identFrom(resolvedType))
						tmp += fmt.Sprintf("on%s(e.value.(%s))\nreturn\n", identFrom(resolvedType), gtype)
						tmp += fmt.Sprintf("}\n")
					}
					tmp += "}\n"

					blk.AddStatements(wdl.RawStmt(tmp))
					blk.AddStatements(wdl.RawStmt("\nif _onDefault != nil {\n_onDefault(e.value)\n}\n"))
				}))
			}),
		)

		for idx, resolvedType := range union.Types() {
			ord := idx + 1
			strct.AddMethods(
				wdl.NewFunc(func(fn *wdl.Func) {
					fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
						r.SetName("e")
						r.SetTypeDef(strct.AsResolvedType())
					}))
					fn.SetVisibility(wdl.Public)
					fn.SetName("as" + identFrom(resolvedType))
					fn.AddResults(
						wdl.NewParam(func(param *wdl.Param) {
							param.SetTypeDef(resolvedType)
						}),
						wdl.NewParam(func(param *wdl.Param) {
							param.SetTypeDef(e.prog.MustResolveSimple("std", "bool"))
						}),
					)
					fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {
						tmp := golang.NewRFile(golang.NewRenderer(golang.Options{}), union.Pkg().Qualifier())
						gtype := tmp.GoType(resolvedType)
						blk.AddStatements(wdl.RawStmt(fmt.Sprintf("var zero %s\nif e.ordinal==%d {\nreturn e.value.(%s), true}\n\n return zero, false\n", gtype, ord, gtype)))
					}))
				}),
				wdl.NewFunc(func(fn *wdl.Func) {
					fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
						r.SetName("e")
						r.SetTypeDef(strct.AsResolvedType())
					}))
					fn.SetVisibility(wdl.Public)
					fn.SetName("with" + identFrom(resolvedType))
					fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
						param.SetTypeDef(resolvedType)
						param.SetName("v")
					}))
					fn.AddResults(
						wdl.NewParam(func(param *wdl.Param) {
							param.SetTypeDef(strct.AsResolvedType())
						}),
					)
					fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {
						blk.AddStatements(wdl.RawStmt(fmt.Sprintf("e.ordinal=%d\ne.value=v\nreturn e", ord)))
					}))
				}),
			)
		}
	})

	union.Pkg().AddFiles(wdl.NewFile(func(file *wdl.File) {
		file.AddImport("json", "encoding/json")
		file.AddImport("fmt", "fmt")

		file.SetName(strings.ToLower(union.Name().String()) + ".gen.go")
		file.SetPath(union.File().Path())
		file.SetModified(true)
		file.SetPreamble(wdl.NewComment(func(comment *wdl.Comment) {
			comment.AddLines(wdl.NewCommentLine(func(line *wdl.CommentLine) {
				line.SetText(e.preamble)
			}))
		}))
		file.AddTypeDefs(uStruct)

		// free but generic match function
		file.AddTypeDefs(wdl.NewFunc(func(fn *wdl.Func) {
			fn.SetName("Match" + uStruct.Name())
			fn.SetVisibility(wdl.Public)
			fn.AddTypeParams(wdl.NewParam(func(param *wdl.Param) {
				param.SetName("R")
				param.SetTypeDef(e.prog.MustResolveSimple("std", "any"))
			}))
			fn.AddResults(wdl.NewParam(func(param *wdl.Param) {
				param.SetName("R")
			}))

			fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
				param.SetName("e")
				param.SetTypeDef(uStruct.AsResolvedType())
			}))

			for _, resolvedType := range union.Types() {
				fn.AddArgs(
					wdl.NewParam(func(param *wdl.Param) {
						param.SetName("on" + identFrom(resolvedType))
						param.SetTypeDef(wdl.NewFunc(func(fn *wdl.Func) {
							fn.SetPkg(union.Pkg())
							fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
								param.SetTypeDef(resolvedType)
							}))
							fn.AddResults(wdl.NewParam(func(param *wdl.Param) {
								param.SetName("R")
							}))
						}).AsResolvedType())
					}),
				)
			}

			fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
				param.SetName("_onDefault") // avoid accidental name collisions using prefix
				param.SetTypeDef(wdl.NewFunc(func(fn *wdl.Func) {
					fn.SetPkg(union.Pkg())
					fn.AddArgs(wdl.NewParam(func(param *wdl.Param) {
						param.SetTypeDef(e.prog.MustResolveSimple("std", "any"))
					}))
					fn.AddResults(wdl.NewParam(func(param *wdl.Param) {
						param.SetName("R")
					}))
				}).AsResolvedType())
			}))

			fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {
				tmp := "switch e.ordinal {\n"
				for idx, resolvedType := range union.Types() {
					ord := idx + 1
					rtmp := golang.NewRFile(golang.NewRenderer(golang.Options{}), union.Pkg().Qualifier())
					gtype := rtmp.GoType(resolvedType)

					tmp += fmt.Sprintf("case %d:\n", ord)
					tmp += fmt.Sprintf("if on%s !=nil {\n", identFrom(resolvedType))
					tmp += fmt.Sprintf("return on%s(e.value.(%s))\n", identFrom(resolvedType), gtype)
					tmp += fmt.Sprintf("}\n")
				}
				tmp += "}\n"

				blk.AddStatements(wdl.RawStmt("if _onDefault == nil{\npanic(`missing default match: cannot guarantee exhaustive matching`)\n}\n\n"))
				blk.AddStatements(wdl.RawStmt(tmp))
				blk.AddStatements(wdl.RawStmt("\nreturn _onDefault(e.value)\n"))
			}))
		}))
	}))

	switch opts.JSONRepresentation {
	case internallyTagged:
		e.goTaggedUnionJSONInternallyTagged(union, uStruct, opts.TagName)
	default:
		return fmt.Errorf("no such json tag variant supported: %s", opts.JSONRepresentation)
	}

	return nil
}

// see also https://serde.rs/enum-representations.html#internally-tagged
func (e *Engine) goTaggedUnionJSONInternallyTagged(union *wdl.Union, uStruct *wdl.Struct, tagAttrName string) {
	uStruct.AddMethods(
		wdl.NewFunc(func(fn *wdl.Func) {
			fn.SetName("MarshalJSON")
			fn.SetVisibility(wdl.Public)
			fn.AddResults(
				wdl.NewParam(func(param *wdl.Param) {
					slice := e.prog.MustResolveSimple("std", "Slice")
					slice.AddParams(e.prog.MustResolveSimple("std", "byte"))
					param.SetTypeDef(slice)
				}),
				wdl.NewParam(func(param *wdl.Param) {
					param.SetTypeDef(e.prog.MustResolveSimple("std", "error"))
				}),
			)
			fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
				r.SetName("e")
				r.SetTypeDef(uStruct.AsResolvedType())
			}))
			fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {
				blk.AddStatements(
					wdl.RawStmt("if e.ordinal == 0 {\nreturn nil, fmt.Errorf(\"marshalling a zero value is not allowed\")\n}\n\n"),
					wdl.RawStmt("return json.Marshal(e)\n"),
				)

			}))
		}),

		wdl.NewFunc(func(fn *wdl.Func) {
			fn.SetName("UnmarshalJSON")
			fn.SetVisibility(wdl.Public)
			fn.AddArgs(
				wdl.NewParam(func(param *wdl.Param) {
					slice := e.prog.MustResolveSimple("std", "Slice")
					slice.AddParams(e.prog.MustResolveSimple("std", "byte"))
					param.SetTypeDef(slice)
					param.SetName("bytes")
				}),
			)
			fn.AddResults(
				wdl.NewParam(func(param *wdl.Param) {
					param.SetTypeDef(e.prog.MustResolveSimple("std", "error"))
				}),
			)
			fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
				r.SetName("e")
				r.SetTypeDef(uStruct.AsResolvedType())
				r.TypeDef().SetPointer(true)
			}))
			fn.SetBody(wdl.NewBlock(func(blk *wdl.Block) {
				blk.AddStatements(
					wdl.RawStmt("typeOnly := struct {\n\t\tType string `json:\""+tagAttrName+"\"`\n\n}{}\n"),
					wdl.RawStmt(`if err := json.Unmarshal(bytes, &typeOnly); err != nil {
		return err
	}`),
				)

				tmp := "\nswitch typeOnly.Type {\n"
				for idx, resolvedType := range union.Types() {
					ord := idx + 1
					rtmp := golang.NewRFile(golang.NewRenderer(golang.Options{}), union.Pkg().Qualifier())
					gtype := rtmp.GoType(resolvedType)

					tmp += fmt.Sprintf("case \"%s\":\n", identFrom(resolvedType))
					tmp += fmt.Sprintf("var value %s\n", gtype)
					tmp += fmt.Sprintf("if err:=json.Unmarshal(bytes, &value);err !=nil {\n")
					tmp += fmt.Sprintf("return fmt.Errorf(\"cannot unmarshal variant '%s'\")\n", gtype)
					tmp += fmt.Sprintf("}\n")
					tmp += fmt.Sprintf("e.ordinal=%d\n", ord)
				}
				tmp += "default:\nreturn fmt.Errorf(\"unknown type variant name '%s'\",typeOnly.Type)"
				tmp += "}\n\nreturn nil\n"

				blk.AddStatements(wdl.RawStmt(tmp))
			}))
		}),
	)
}

func identFrom(resolvedType *wdl.ResolvedType) wdl.Identifier {
	if len(resolvedType.Params()) == 0 {
		return resolvedType.Name()
	}

	var compoundName wdl.Identifier
	for _, r := range resolvedType.Params() {
		compoundName += r.Name()
	}
	return compoundName + resolvedType.Name()
}
