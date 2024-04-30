package golang

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
)

func (r *RFile) renderFunc(f *wdl.Func, w *render.Writer) error {
	r.parent.writeComments(w, f.Comment())
	w.Printf("func ")
	if f.Receiver() != nil {
		w.Printf("(%s %s) ", f.Receiver().Name(), r.GoType(f.Receiver().TypeDef()))
	}
	// func name
	w.Print(goAccessorName(f))

	// func params
	w.Printf("(")
	for _, param := range f.Args() {
		if param.Name() != "" {
			w.Printf("%s ", param.Name())
		}
		w.Print(r.GoType(param.TypeDef()))
		w.Print(",")
	}
	w.Printf(")")

	// func result
	if len(f.Results()) > 0 {
		w.Printf("(")
		for _, param := range f.Results() {
			if param.Name() != "" {
				w.Printf("%s ", param.Name())
			}
			w.Print(r.GoType(param.TypeDef()))
			w.Print(",")
		}
		w.Printf(")")
	}

	// func block body
	if f.Body() != nil {
		w.Printf("{\n")
		for _, statement := range f.Body().Statements() {
			if err := r.renderStmt(statement, w); err != nil {
				return err
			}
		}
		w.Printf("}\n")
	}

	return nil
}
