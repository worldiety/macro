package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"strings"
)

// formatComment replaces a '...' prefix with the ellipsisName and prefixes all lines
// with a '// '. If doc is empty, the empty string is returned.
func formatComment(ellipsisName, doc string, indent int) string {
	doc = strings.TrimSpace(doc)
	if len(doc) > 0 {
		tmp := &strings.Builder{}
		if strings.HasPrefix(doc, "...") {
			tmp.WriteString(ellipsisName)
			tmp.WriteString(" ")
			tmp.WriteString(strings.TrimSpace(doc[3:]))
		} else {
			tmp.WriteString(doc)
		}
		str := tmp.String()
		tmp.Reset()
		for i := 0; i < indent; i++ {
			tmp.WriteString(" ")
		}
		tmp.WriteString("/**\n") // this is JSDoc style, see also https://google.github.io/styleguide/tsguide.html#jsdoc-general-form
		for _, line := range strings.Split(str, "\n") {
			for i := 0; i < indent; i++ {
				tmp.WriteString(" ")
			}
			tmp.WriteString(" * ")
			tmp.WriteString(line)
			tmp.WriteString("\n")
		}
		for i := 0; i < indent; i++ {
			tmp.WriteString(" ")
		}
		tmp.WriteString(" */\n")

		return tmp.String()
	}

	return ""
}

// DeEllipsis replaces a ... with the according text.
func DeEllipsis(ellipsisName, doc string) string {
	tmp := &strings.Builder{}
	if strings.HasPrefix(doc, "...") {
		tmp.WriteString(ellipsisName)
		tmp.WriteString(" ")
		tmp.WriteString(strings.TrimSpace(doc[3:]))
	} else {
		tmp.WriteString(doc)
	}

	return tmp.String()
}

func (r *Renderer) writeCommentNode(w *render.Writer, isPkg bool, name string, intent int, comment *wdl.Comment) {
	if comment == nil {
		return
	}

	r.writeComment(w, isPkg, name, comment.String(), intent)
}

func (r *Renderer) writeComments(w *render.Writer, intent int, lines []*wdl.CommentLine) {
	if lines == nil {
		return
	}

	r.writeCommentNode(w, false, "", intent, wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(lines)
	}))
}

func (r *Renderer) writeComment(w *render.Writer, isPkg bool, name, doc string, intent int) {
	if isPkg {
		name = "Package " + name
	}

	myDoc := formatComment(name, doc, intent)
	if doc != "" {
		w.Printf(myDoc)
	}
}
