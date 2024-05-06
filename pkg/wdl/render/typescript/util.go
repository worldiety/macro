package typescript

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"path/filepath"
	"unicode"
)

func (r *RFile) TsType(rtype *wdl.ResolvedType) string {
	if rtype == nil {
		return "'no type resolved'"
	}
	// TODO currently we anticipate that all typescript is generated in the same flat directory

	if rtype.TypeParam() {
		return rtype.Name().String()
	}

	r.Use(rtype)
	switch def := rtype.TypeDef().(type) {
	case *wdl.BaseType:
		switch def.Kind() {
		case wdl.TString:
			return "string"
		case wdl.TInt:
			fallthrough
		case wdl.TInt64:
			return "number"
		case wdl.TAny:
			return "any"
		case wdl.TBool:
			return "bool"
		case wdl.TByte:
			return "byte"
		default:
			panic(fmt.Errorf("implement me: %v", def.Kind()))
		}
	case *wdl.Func:
		tmp := &render.Writer{}
		if err := r.renderFunc(def, tmp); err != nil {
			panic(err) // TODO ???
		}
		return tmp.String()
	case *wdl.TypeParam:
		return def.Name().String()

	default:
		r.AddImport(rtype.Name(), wdl.PkgImportQualifier(filepath.Join(string(r.selfImportPath), tsLowerNameStr(rtype.Name().String()))))

		if rtype.Pkg().Name() == "std" {
			switch rtype.Name() {
			case "error":
				return "error"
			case "Slice":
				if len(rtype.Params()) != 1 {
					panic(fmt.Errorf("invalid Slice type param: %#v", rtype))
				}
				return "[]" + r.TsType(rtype.Params()[0])
			}

		}

		tmp := rtype.Name().String()
		if len(rtype.Params()) > 0 {
			tmp += "<"
			for _, resolvedType := range rtype.Params() {
				tmp += r.TsType(resolvedType)
				tmp += ","
			}
			tmp += ">"
		}
		return tmp
	}

}

func tsUpperName(f interface {
	Visibility() wdl.Visibility
	Name() wdl.Identifier
}) string {
	if f.Name() == "" {
		return ""
	}

	return tsUpperNameStr(string(f.Name()))
}

func tsUpperNameStr(s string) string {
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

func tsLowerNameStr(s string) string {
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func tsLowerName(f interface {
	Visibility() wdl.Visibility
	Name() wdl.Identifier
}) string {
	if f.Name() == "" {
		return ""
	}

	return tsLowerNameStr(string(f.Name()))
}
