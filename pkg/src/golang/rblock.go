package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

// renderBlock emits a block and all contained nodes.
func (r *Renderer) renderBlock(node *ast.Block, w *render.BufferedWriter) error {
	r.writeCommentNode(w, false, "", node.ObjComment)
	w.Printf("{\n")
	for _, n := range node.Nodes {
		if err := r.renderNode(n, w); err != nil {
			return fmt.Errorf("unable to render node in block: %w", err)
		}
	}

	w.Printf("}\n")

	return nil
}
