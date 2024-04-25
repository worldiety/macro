package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderMacro emits something which is usually evaluated here.
func (r *Renderer) renderMacro(node *ast.Macro, w *render.BufferedWriter) error {
	r.writeCommentNode(w, false, "", node.Comment())
	if node.Func != nil {
		actualNodes := node.Children()
		for _, actualNode := range actualNodes {
			if err := r.renderNode(actualNode, w); err != nil {
				return fmt.Errorf("unable to render dynamic macro node: %w", err)
			}
		}
	}

	return nil
}
