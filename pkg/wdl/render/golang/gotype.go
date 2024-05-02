package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"unicode"
)

func (r *RFile) GoType(rtype *wdl.ResolvedType) string {
	var ptr string
	if rtype.Pointer() {
		ptr += "*"
	}

	r.Use(rtype)
	switch def := rtype.TypeDef().(type) {
	case *wdl.BaseType:
		switch def.Kind() {
		case wdl.TString:
			return ptr + "string"
		case wdl.TInt:
			return ptr + "int"
		case wdl.TAny:
			return ptr + "any"
		case wdl.TBool:
			return ptr + "bool"
		case wdl.TByte:
			return ptr + "byte"
		default:
			panic(fmt.Errorf("implement me: %v", def.Kind()))
		}
	case *wdl.Func:
		tmp := &render.Writer{}
		if err := r.renderFunc(def, tmp); err != nil {
			panic(err) // TODO ???
		}
		return ptr + tmp.String()
	default:
		if r.selfImportPath == rtype.Pkg().Qualifier() {
			// just a package local type
			return ptr + rtype.Name().String()
		}

		if rtype.Pkg().Name() == "std" {
			switch rtype.Name() {
			case "error":
				return ptr + "error"
			case "Slice":
				if len(rtype.Params()) != 1 {
					panic(fmt.Errorf("invalid Slice type param: %#v", rtype))
				}
				return ptr + "[]" + r.GoType(rtype.Params()[0])
			}

		}

		tmp := ptr + rtype.Pkg().Name().String() + "." + rtype.Name().String()
		if len(rtype.Params()) > 0 {
			tmp += "["
			for _, resolvedType := range rtype.Params() {
				tmp += r.GoType(resolvedType)
				tmp += ","
			}
			tmp += "]"
		}
		return tmp
	}

}

func goAccessorName(f interface {
	Visibility() wdl.Visibility
	Name() wdl.Identifier
}) string {
	if f.Name() == "" {
		return ""
	}

	switch f.Visibility() {
	case wdl.Public:
		return string(unicode.ToUpper(rune(f.Name().String()[0]))) + f.Name().String()[1:]
	default:
		return string(unicode.ToLower(rune(f.Name().String()[0]))) + f.Name().String()[1:]
	}
}
