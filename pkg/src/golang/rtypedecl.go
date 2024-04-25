package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
	"github.com/worldiety/macro/pkg/src/stdlib"
	"reflect"
)

func (r *Renderer) renderTypeDecl(node ast.TypeDecl, w *render.BufferedWriter) error {
	importer := r.importer(node)

	switch t := node.(type) {
	case *ast.SimpleTypeDecl:
		w.Printf(string(importer.shortify(fromStdlib(t.SimpleName))))
	case *ast.TypeDeclPtr:
		w.Printf("*")
		if err := r.renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
	case *ast.SliceTypeDecl:
		w.Printf("[]")
		if err := r.renderTypeDecl(t.TypeDecl, w); err != nil {
			return err
		}
	case *ast.GenericTypeDecl:
		if err := r.renderTypeDecl(t.TypeDecl, w); err != nil {
			return err
		}

		builtInHandled := false
		if std, ok := t.TypeDecl.(*ast.SimpleTypeDecl); ok {
			switch std.Name() {
			case stdlib.Map:
				w.Printf("[")
				if err := r.renderTypeDecl(t.Params()[0], w); err != nil {
					return err
				}

				w.Printf("]")
				if err := r.renderTypeDecl(t.Params()[1], w); err != nil {
					return err
				}

				builtInHandled = true
			case stdlib.List:
				if err := r.renderTypeDecl(t.Params()[0], w); err != nil {
					return err
				}

				builtInHandled = true
			}

		}

		if !builtInHandled {
			w.Printf("[")
			for i, decl := range t.Params() {
				if err := r.renderTypeDecl(decl, w); err != nil {
					return err
				}
				if i < len(t.Params())-1 {
					w.Printf(",")
				}
			}
			w.Printf("]")
		}

	case *ast.ChanTypeDecl:
		w.Printf("chan ")
		if err := r.renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}

	case *ast.ArrayTypeDecl:
		w.Printf("[%d]", t.ArrayLen)
		if err := r.renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}

		/*	case *ast.FuncTypeDecl:
			w.Printf("func ")
			if err := r.renderFunc(t.Func(), w); err != nil {
				return err
			}*/
	default:
		panic("not yet implemented: " + reflect.TypeOf(t).String())
	}

	return nil
}
