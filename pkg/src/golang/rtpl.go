package golang

import (
	"bytes"
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
	"text/template"
)

// renderTpl executes and emits the template text.
func (r *Renderer) renderTpl(node *ast.Tpl, w *render.BufferedWriter) error {
	importer := r.importer(node)
	ctx := &tplRenderContext{
		importer: importer,
		tpl:      node,
	}

	tmpl, err := template.New(node.ObjPos.String()).Parse(node.Template)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, ctx); err != nil {
		return fmt.Errorf("cannot execute template: %w", err)
	}

	w.Print(buf.String())

	return nil
}

// ensure that we always implement the full contract
var _ ast.TplContext = (*tplRenderContext)(nil)

type tplRenderContext struct {
	importer *importer
	tpl      *ast.Tpl
}

func (t *tplRenderContext) Get(key string) interface{} {
	return t.tpl.Values[key]
}

func (t *tplRenderContext) Use(name string) string {
	goType := fromStdlib(ast.Name(name))
	return string(t.importer.shortify(goType))
}

func (t *tplRenderContext) Self() *ast.Tpl {
	return t.tpl
}
