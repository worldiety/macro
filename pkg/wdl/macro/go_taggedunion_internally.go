package macro

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render/golang"
	"strconv"
)

// see also https://serde.rs/enum-representations.html#internally-tagged
func (m *GoTaggedUnion) goTaggedUnionJSONInternallyTagged(opts goTaggedUnionParams, union *wdl.Union, uStruct *wdl.Struct, tagAttrName string) {
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
					wdl.RawStmt("buf,err:= json.Marshal(e.value)\n\tif err !=nil{\n\t\treturn nil,err\n\t}\nvar prefix []byte\n"),
				)
				tmp := "\nswitch e.ordinal {\n"
				for idx, resolvedType := range union.Types() {
					ord := idx + 1
					strCaseConst := identFrom(resolvedType)
					if len(opts.Names) > 0 {
						strCaseConst = wdl.Identifier(opts.Names[idx])
					}
					tmp += fmt.Sprintf("case %d:\n", ord)
					tmp += fmt.Sprintf("prefix = []byte(`{\"type\":%s`)\n", strconv.Quote(strCaseConst.String()))
				}
				tmp += "}\n"
				blk.Add(wdl.RawStmt(tmp))
				blk.Add(wdl.RawStmt(`

if len(buf)>2{
	// we expect an empty object like {} or at least an object with a single attribute, which requires a , separator
	prefix=append(prefix,',')
}
buf= append(buf[1:], prefix...)
	copy(buf[len(prefix):], buf)
	copy(buf,prefix)

	return buf,nil`))
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

					strCaseConst := identFrom(resolvedType)
					if len(opts.Names) > 0 {
						strCaseConst = wdl.Identifier(opts.Names[idx])
					}
					tmp += fmt.Sprintf("case \"%s\":\n", strCaseConst)
					tmp += fmt.Sprintf("var value %s\n", gtype)
					tmp += fmt.Sprintf("if err:=json.Unmarshal(bytes, &value);err !=nil {\n")
					tmp += fmt.Sprintf("return fmt.Errorf(\"cannot unmarshal variant '%s'\")\n", gtype)
					tmp += fmt.Sprintf("}\n")
					tmp += fmt.Sprintf("e.ordinal=%d\n", ord)
					tmp += fmt.Sprintf("e.value=value\n")
				}
				tmp += "default:\nreturn fmt.Errorf(\"unknown type variant name '%s'\",typeOnly.Type)"
				tmp += "}\n\nreturn nil\n"

				blk.Add(wdl.RawStmt(tmp))
			}))
		}),
	)
}
