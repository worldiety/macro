package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderSym emits an imported qualifier.
func (r *Renderer) renderCompLit(node *ast.CompLit, w *render.BufferedWriter) error {
	if node.Type != nil {
		if err := r.renderNode(node.Type, w); err != nil {
			return fmt.Errorf("unable to render type: %w", err)
		}
	}

	w.Print("{")
	for i, element := range node.Elements {
		if err := r.renderNode(element, w); err != nil {
			return fmt.Errorf("unable to render composite elem: %w", err)
		}

		if i < len(node.Elements)-1 {
			w.Print(",")
		}
	}
	w.Print("}")
	return nil
}
