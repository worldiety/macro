package macro

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render/golang"
	"strconv"
)

// see also https://serde.rs/enum-representations.html#adjacently-tagged
func (m *GoTaggedUnion) goTaggedUnionJSONAdjacentlyTagged(opts goTaggedUnionParams, union *wdl.Union, uStruct *wdl.Struct) {
	uStruct.AddMethods(
		wdl.NewFunc(func(fn *wdl.Func) {
			fn.SetName("MarshalJSON")
			fn.SetVisibility(wdl.Public)
			fn.AddResults(
				wdl.NewParam(func(param *wdl.Param) {
					slice := m.prog.MustResolveSimple("std", "Slice")
					slice.AddParams(m.prog.MustResolveSimple("std", "byte"))
					param.SetTypeDef(slice)
				}),
				wdl.NewParam(func(param *wdl.Param) {
					param.SetTypeDef(m.prog.MustResolveSimple("std", "error"))
				}),
			)
			fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
				r.SetName("e")
				r.SetTypeDef(uStruct.AsResolvedType())
			}))
			fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {
				blk.Add(
					wdl.RawStmt("if e.ordinal == 0 {\nreturn nil, fmt.Errorf(\"marshalling a zero value is not allowed\")\n}\n\n"),
					wdl.RawStmt("// note, that by definition, this kind of encoding does not work with union types which evaluates to null, arrays or primitives.\n// Chose adjacent encoding instead.\n"),
					wdl.RawStmt(fmt.Sprintf("type adjacentlyTagged[T any] struct {\n\tType  string `json:%s`\n\tValue T      `json:%s`\n}", strconv.Quote(opts.TagName), strconv.Quote(opts.Content))),
					wdl.RawStmt("\n"),
				)
				tmp := "\nswitch e.ordinal {\n"
				for idx, resolvedType := range union.Types() {
					ord := idx + 1
					strCaseConst := identFrom(resolvedType)
					if len(opts.Names) > 0 {
						strCaseConst = wdl.Identifier(opts.Names[idx])
					}

					rtmp := golang.NewRFile(golang.NewRenderer(golang.Options{}), union.Pkg().Qualifier())
					gtype := rtmp.GoType(resolvedType)

					tmp += fmt.Sprintf("case %d:\n", ord)
					tmp += fmt.Sprintf("return json.Marshal(adjacentlyTagged[%s]{\n\t\t\tType:  %s,\n\t\t\tValue: e.value.(%s),\n\t\t})\n", gtype, strconv.Quote(strCaseConst.String()), gtype)
				}
				tmp += "default:\nreturn nil,fmt.Errorf(\"unknown type ordinal variant '%d'\",e.ordinal)"

				tmp += "}\n"
				blk.Add(wdl.RawStmt(tmp))
			}))
		}),

		wdl.NewFunc(func(fn *wdl.Func) {
			fn.SetName("UnmarshalJSON")
			fn.SetVisibility(wdl.Public)
			fn.AddArgs(
				wdl.NewParam(func(param *wdl.Param) {
					slice := m.prog.MustResolveSimple("std", "Slice")
					slice.AddParams(m.prog.MustResolveSimple("std", "byte"))
					param.SetTypeDef(slice)
					param.SetName("bytes")
				}),
			)
			fn.AddResults(
				wdl.NewParam(func(param *wdl.Param) {
					param.SetTypeDef(m.prog.MustResolveSimple("std", "error"))
				}),
			)
			fn.SetReceiver(wdl.NewParam(func(r *wdl.Param) {
				r.SetName("e")
				r.SetTypeDef(uStruct.AsResolvedType())
				r.TypeDef().SetPointer(true)
			}))
			fn.SetBody(wdl.NewBlockStmt(func(blk *wdl.BlockStmt) {
				blk.Add(
					wdl.RawStmt("typeOnly := struct {\n\t\tType string `json:\""+opts.TagName+"\"`\n\n}{}\n"),
					wdl.RawStmt(`if err := json.Unmarshal(bytes, &typeOnly); err != nil {
		return err
	}`),
					wdl.RawStmt(fmt.Sprintf("\ntype adjacentlyTagged[T any] struct {\n\tType  string `json:%s`\n\tValue T      `json:%s`\n}", strconv.Quote(opts.TagName), strconv.Quote(opts.Content))),
				)

				tmp := "\nswitch typeOnly.Type {\n"
				for idx, resolvedType := range union.Types() {
					ord := idx + 1
					rtmp := golang.NewRFile(golang.NewRenderer(golang.Options{}), union.Pkg().Qualifier())
					gtype := rtmp.GoType(resolvedType)

					strCaseConst := identFrom(resolvedType)
					if len(opts.Names) > 0 {
						strCaseConst = wdl.Identifier(opts.Names[idx])
					}
					tmp += fmt.Sprintf("case \"%s\":\n", strCaseConst)
					tmp += fmt.Sprintf("var value adjacentlyTagged[%s]\n", gtype)
					tmp += fmt.Sprintf("if err:=json.Unmarshal(bytes, &value);err !=nil {\n")
					tmp += fmt.Sprintf("return fmt.Errorf(\"cannot unmarshal variant '%s'\")\n", gtype)
					tmp += fmt.Sprintf("}\n")
					tmp += fmt.Sprintf("e.ordinal=%d\n", ord)
					tmp += fmt.Sprintf("e.value=value.Value\n")
				}
				tmp += "default:\nreturn fmt.Errorf(\"unknown type variant name '%s'\",typeOnly.Type)"
				tmp += "}\n\nreturn nil\n"

				blk.Add(wdl.RawStmt(tmp))
			}))
		}),
	)
}
