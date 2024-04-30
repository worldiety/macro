package macro

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render/golang"
	"strings"
)

func (e *Engine) goTaggedUnion(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {
	if _, ok := def.(*wdl.Interface); ok {
		// TODO this case happens, cannot decide properly if that is correct or wrong
		return nil
	}
	union, ok := def.(*wdl.Union)
	if !ok {
		return fmt.Errorf("expected union definition")
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
							param.SetName("on" + resolvedType.Name())
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
						tmp += fmt.Sprintf("if on%s !=nil {\n", resolvedType.Name())
						tmp += fmt.Sprintf("on%s(e.value.(%s))\nreturn\n", resolvedType.Name(), gtype)
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
					fn.SetName("as" + resolvedType.Name())
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
					fn.SetName("with" + resolvedType.Name())
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
		file.SetName(strings.ToLower(union.Name().String()) + ".gen.go")
		file.SetPath(union.File().Path())
		file.SetModified(true)
		file.SetPreamble(wdl.NewComment(func(comment *wdl.Comment) {
			comment.AddLines(wdl.NewCommentLine(func(line *wdl.CommentLine) {
				line.SetText(e.preamble)
			}))
		}))
		file.AddTypeDefs(uStruct)

	}))

	return nil
}
