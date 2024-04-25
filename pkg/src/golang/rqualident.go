package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderQualIdent emits an imported qualifier.
func (r *Renderer) renderQualIdent(node *ast.QualIdent, w *render.BufferedWriter) error {
	importer := r.importer(node)
	renamedQualifier := importer.shortify(fromStdlib(ast.Name(node.Qualifier) + "._")).Qualifier()
	w.Printf(renamedQualifier)

	return nil
}
