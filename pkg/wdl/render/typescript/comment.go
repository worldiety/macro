package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"strings"
)

// formatComment replaces a '...' prefix with the ellipsisName and prefixes all lines
// with a '// '. If doc is empty, the empty string is returned.
func formatComment(ellipsisName, doc string) string {
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
		for _, line := range strings.Split(str, "\n") {
			tmp.WriteString("// ")
			tmp.WriteString(line)
			tmp.WriteString("\n")
		}

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

func (r *Renderer) writeCommentNode(w *render.Writer, isPkg bool, name string, comment *wdl.Comment) {
	if comment == nil {
		return
	}

	r.writeComment(w, isPkg, name, comment.String())
}

func (r *Renderer) writeComments(w *render.Writer, lines []*wdl.CommentLine) {
	if lines == nil {
		return
	}

	r.writeCommentNode(w, false, "", wdl.NewComment(func(comment *wdl.Comment) {
		comment.SetLines(lines)
	}))
}

func (r *Renderer) writeComment(w *render.Writer, isPkg bool, name, doc string) {
	if isPkg {
		name = "Package " + name
	}

	myDoc := formatComment(name, doc)
	if doc != "" {
		w.Printf(myDoc)
	}
}
