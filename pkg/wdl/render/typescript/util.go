package typescript

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"log/slog"
	"path/filepath"
	"strings"
	"unicode"
)

func looksLikeOptionHack(rtype *wdl.ResolvedType) bool {
	return rtype.Name() == "Option" && len(rtype.Params()) == 1
}

func (r *RFile) TsType(rtype *wdl.ResolvedType) string {
	if rtype == nil {
		return "'no type resolved'"
	}
	// TODO currently we anticipate that all typescript is generated in the same flat directory

	if rtype.TypeParam() {
		return rtype.Name().String()
	}

	// special Option hack, which just applies to "all Option[T]" types
	if looksLikeOptionHack(rtype) {
		return r.TsType(rtype.Params()[0])
	}

	r.Use(rtype)
	switch def := rtype.TypeDef().(type) {
	case *wdl.BaseType:
		switch def.Kind() {
		case wdl.TString:
			return "string"
		case wdl.TByte:

			return "number /*uint8*/"
		case wdl.TInt:
			return "number /*int*/"
		case wdl.TInt64:
			return "number /*int64*/"
		case wdl.TAny:
			return "any"
		case wdl.TBool:
			return "boolean"

		case wdl.TFloat32:
			return "number /*float32*/"
		case wdl.TFloat64:
			return "number /*float64*/"
		default:
			panic(fmt.Errorf("implement me: %v", def.Kind()))
		}
	case *wdl.Func:
		//tmp := &render.Writer{}
		//if err := r.renderFunc(def, tmp); err != nil {
		//	panic(err) // TODO ???
		//}
		//return tmp.String()
		slog.Info("ignored typescript func emitter") // TODO we have some recursion problems, probably type and method name collisions in NAGO/ora protocol
		return ""
	case *wdl.TypeParam:
		return def.Name().String()

	default:

		if rtype.Pkg().Name() == "std" {
			switch rtype.Name() {
			case "error":
				return "error"
			case "Slice":
				if len(rtype.Params()) != 1 {
					panic(fmt.Errorf("invalid Slice type param: %#v", rtype))
				}
				return r.TsType(rtype.Params()[0]) + "[]"
			case "Map":
				if len(rtype.Params()) != 2 {
					panic(fmt.Errorf("invalid map type param: %#v", rtype))
				}
				return "Record<" + r.TsType(rtype.Params()[0]) + "," + r.TsType(rtype.Params()[1]) + ">"
			}

		}

		r.AddImport(rtype.Name(), wdl.PkgImportQualifier(filepath.Join(string(r.selfImportPath), tsLowerNameStr(rtype.Name().String()))))

		tmp := rtype.Name().String()
		if len(rtype.Params()) > 0 {
			tmp += "<"
			for i, resolvedType := range rtype.Params() {
				tmp += r.TsType(resolvedType)
				if i != len(rtype.Params())-1 {
					tmp += ", "
				}
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
	name := f.Name().String()
	if strings.HasPrefix(name, "_") {
		name = name[1:]
	}

	return tsUpperNameStr(name)
}

func tsUpperNameStr(s string) string {
	first, rest := wdl.SplitFirstRune(s)
	return string(unicode.ToUpper(first)) + rest
}

func tsLowerNameStr(s string) string {
	s = strings.TrimLeft(s, "._")
	first, rest := wdl.SplitFirstRune(s)
	return string(unicode.ToLower(first)) + rest
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
