package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"unicode"
)

func (r *RFile) renderStruct(def *wdl.Struct, w *render.Writer) error {

	r.parent.writeCommentNode(w, false, "", wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(def.Comment())
	}))
	w.Printf("type %s struct {\n", goAccessorName(def))
	for _, field := range def.Fields() {
		w.Printf("%s %s\n", goAccessorName(field), r.GoType(field.TypeDef()))
	}
	w.Printf("}\n\n")

	for _, f := range def.Methods() {
		if err := r.renderFunc(f, w); err != nil {
			return err
		}
		w.Print("\n")
	}

	return nil
}

func (r *RFile) GoType(rtype *wdl.ResolvedType) string {
	r.Use(rtype)
	switch def := rtype.TypeDef().(type) {
	case *wdl.BaseType:
		switch def.Kind() {
		case wdl.TString:
			return "string"
		case wdl.TInt:
			return "int"
		case wdl.TAny:
			return "any"
		case wdl.TBool:
			return "bool"
		default:
			panic(fmt.Errorf("implement me: %v", def.Kind()))
		}
	case *wdl.Func:
		tmp := &render.Writer{}
		if err := r.renderFunc(def, tmp); err != nil {
			panic(err) // TODO ???
		}
		return tmp.String()
	default:
		if r.selfImportPath == rtype.Pkg().Qualifier() {
			// just a package local type
			return rtype.Name().String()
		}

		return rtype.Pkg().Name().String() + "." + rtype.Name().String()
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
