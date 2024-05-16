package macro

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render/golang"
	"strings"
)

type jsonRepr string

const (
	internallyTagged jsonRepr = "internal"
	adjacentTagged   jsonRepr = "adjacent"
)

type goTaggedUnionParams struct {
	JSONRepresentation jsonRepr `json:"json"`    // either "intern" or "adjacent", default is "adjacent"
	TagName            string   `json:"tag"`     // used by internally and adjacently tagged, default is "type"
	Content            string   `json:"content"` // used by adjacently tagged, default is "content"
	Names              []string `json:"names"`
}

type GoTaggedUnion struct {
	prog     *wdl.Program
	preamble string
}

func NewGoTaggedUnion(prog *wdl.Program, preamble string) *GoTaggedUnion {
	return &GoTaggedUnion{prog: prog, preamble: preamble}
}

func (m *GoTaggedUnion) Names() []wdl.MacroName {
	return []wdl.MacroName{"go.TaggedUnion"}
}

func (m *GoTaggedUnion) Expand(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {
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

	if opts.JSONRepresentation == "" {
		opts.JSONRepresentation = adjacentTagged
	}

	if opts.TagName == "" {
		opts.TagName = "type"
	}

	if opts.Content == "" {
		opts.Content = "content"
	}

	if len(opts.Names) > 0 && len(opts.Names) != len(union.Types()) {
		return fmt.Errorf("names and union types have different length: %d names vs %d types", len(opts.Names), len(union.Types()))
	}

	uStruct := wdl.NewStruct(func(strct *wdl.Struct) {
		strct.SetPkg(union.Pkg())
		strct.SetVisibility(union.Visibility())
		strct.SetName(wdl.Identifier(stripName(union.Name().String())))
		strct.SetComment(union.Comment())
		strct.AddFields(
			wdl.NewField(func(field *wdl.Field) {
				field.SetName("ordinal")
				field.SetTypeDef(m.prog.MustResolveSimple("std", "int"))
			}),
			wdl.NewField(func(field *wdl.Field) {
				field.SetName("value")
				field.SetTypeDef(m.prog.MustResolveSimple("std", "any"))
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
						param.SetTypeDef(m.prog.MustResolveSimple("std", "any"))
					}),
				)
				fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {

					blk.Add(wdl.RawStmt(fmt.Sprintf("return e.value")))
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
						param.SetTypeDef(m.prog.MustResolveSimple("std", "int"))
					}),
				)
				fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {

					blk.Add(wdl.RawStmt(fmt.Sprintf("return e.ordinal")))
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
						param.SetTypeDef(m.prog.MustResolveSimple("std", "bool"))
					}),
				)
				fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {

					blk.Add(wdl.RawStmt(fmt.Sprintf("return e.ordinal>0")))
				}))
			}),
			wdl.NewFunc(func(fn *wdl.Func) {
				fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
					r.SetName("e")
					r.SetTypeDef(strct.AsResolvedType())
				}))
				fn.SetVisibility(wdl.Public)
				fn.SetName("Switch")
				fn.SetComment(wdl.NewSimpleComment("Switch provides an exhaustive and type safe closure callback mechanic. Nil callbacks are allowed. Unmatched branches are delegated into a default case."))
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
							param.SetTypeDef(m.prog.MustResolveSimple("std", "any"))
						}))
					}).AsResolvedType())
				}))

				fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {
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

					blk.Add(wdl.RawStmt(tmp))
					blk.Add(wdl.RawStmt("\nif _onDefault != nil {\n_onDefault(e.value)\n}\n"))
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
							param.SetTypeDef(m.prog.MustResolveSimple("std", "bool"))
						}),
					)
					fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {
						tmp := golang.NewRFile(golang.NewRenderer(golang.Options{}), union.Pkg().Qualifier())
						gtype := tmp.GoType(resolvedType)
						blk.Add(wdl.RawStmt(fmt.Sprintf("var zero %s\nif e.ordinal==%d {\nreturn e.value.(%s), true}\n\n return zero, false\n", gtype, ord, gtype)))
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
					fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {
						blk.Add(wdl.RawStmt(fmt.Sprintf("e.ordinal=%d\ne.value=v\nreturn e", ord)))
					}))
				}),
			)
		}
	})

	union.Pkg().AddFiles(wdl.NewFile(func(file *wdl.File) {
		file.AddImport("json", "encoding/json")
		file.AddImport("fmt", "fmt")
		file.SetMimeType(wdl.MimeTypeGo)

		file.SetName(strings.ToLower(union.Name().String()) + ".gen.go")
		file.SetPath(union.File().Path())
		file.SetModified(true)
		file.SetGenerated(true)
		file.SetPreamble(wdl.NewComment(func(comment *wdl.Comment) {
			comment.AddLines(wdl.NewCommentLine(func(line *wdl.CommentLine) {
				line.SetText(m.preamble)
			}))
		}))

		file.AddTypeDefs(wdl.NewDistinctType(func(dType *wdl.DistinctType) {
			dType.SetName("_")
			dType.SetPkg(file.Pkg())
			dType.SetUnderlying(def.AsResolvedType().TypeDef())
			dType.SetComment(wdl.NewSimpleComment(fmt.Sprintf("This variable is declared to let linters know, that [%s] is used at compile time to generate [%s].", def.Name(), uStruct.Name())))
		}))
		file.AddTypeDefs(uStruct)

		// free but generic match function
		file.AddTypeDefs(wdl.NewFunc(func(fn *wdl.Func) {
			fn.SetName("Match" + uStruct.Name())
			fn.SetVisibility(wdl.Public)
			fn.AddTypeParams(wdl.NewResolvedType(func(rType *wdl.ResolvedType) {
				rType.SetName("R")
				rType.SetTypeParam(true)
				rType.SetTypeDef(m.prog.MustResolveSimple("std", "any").TypeDef())
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
						param.SetTypeDef(m.prog.MustResolveSimple("std", "any"))
					}))
					fn.AddResults(wdl.NewParam(func(param *wdl.Param) {
						param.SetName("R")
					}))
				}).AsResolvedType())
			}))

			fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {
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

				blk.Add(wdl.RawStmt("if _onDefault == nil{\npanic(`missing default match: cannot guarantee exhaustive matching`)\n}\n\n"))
				blk.Add(wdl.RawStmt(tmp))
				blk.Add(wdl.RawStmt("\nreturn _onDefault(e.value)\n"))
			}))
		}))
	}))

	switch opts.JSONRepresentation {
	case internallyTagged:
		m.goTaggedUnionJSONInternallyTagged(opts, union, uStruct, opts.TagName)
	case adjacentTagged:
		m.goTaggedUnionJSONAdjacentlyTagged(opts, union, uStruct)
	default:
		return fmt.Errorf("no such json tag variant supported: %s", opts.JSONRepresentation)
	}

	return nil
}

func identFrom(resolvedType *wdl.ResolvedType) wdl.Identifier {
	rName := wdl.Identifier(golang.MakePublic(resolvedType.Name().String()))
	if len(resolvedType.Params()) == 0 {
		return rName
	}

	var compoundName wdl.Identifier
	for _, r := range resolvedType.Params() {
		compoundName += wdl.Identifier(golang.MakePublic(r.Name().String()))
	}
	return compoundName + rName
}
