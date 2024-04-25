package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
	"strconv"
	"strings"
)

// renderFile generates the code for the entire file.
func (r *Renderer) renderFile(file *ast.File) ([]byte, error) {
	w := &render.BufferedWriter{}

	// file license or whatever
	if file.Preamble != nil {
		r.writeComment(w, false, file.Pkg().Name, file.Preamble.Text)
		w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	}

	// actual package comment
	if file.Comment() != nil {
		r.writeComment(w, true, file.Pkg().Name, file.Comment().Text)
	}

	w.Printf("package %s\n", file.Pkg().Name)

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	tmp := &render.BufferedWriter{}
	for _, node := range file.Nodes {
		switch t := node.(type) {
		case *ast.Func:
			if err := r.renderFunc(t, tmp); err != nil {
				return nil, err
			}

		default:
			if err := r.renderNode(t, tmp); err != nil {
				return nil, err
			}
		}
	}

	importer := r.importer(file)
	if len(importer.namedImports) > 0 || len(file.Imports()) > 0 {
		w.Printf("import (\n")
		for namedImport, qualifier := range importer.namedImports {
			w.Printf("  %s %s\n", namedImport, strconv.Quote(qualifier))
		}

		for _, imp := range file.Imports() {
			comment := strings.TrimSpace(formatComment(imp.Ident, imp.CommentText()))
			multiline := strings.LastIndex(comment, "\n") > 0
			if multiline {
				r.writeComment(w, false, "", comment)
			}

			w.Printf("  %s %s", imp.Ident, strconv.Quote(string(imp.Name)))
			if len(comment) > 0 && !multiline {
				w.Printf(comment)
			}

			w.Printf("\n")
		}

		w.Printf(")\n")
	}

	w.Printf(tmp.String())

	return Format(w.Bytes())
}
