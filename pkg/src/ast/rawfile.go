package ast

import (
	"bytes"
	"fmt"
	"text/template"
)

// A RawFile is an uninterpreted blob.
type RawFile struct {
	Name     string
	MimeType string
	Data     func(file *RawFile) ([]byte, error)
	Obj
}

// NewRawFile just creates a simple raw file wrapper.
func NewRawFile(name, mimeType string, data []byte) *RawFile {
	return &RawFile{
		Name:     name,
		MimeType: mimeType,
		Data: func(_ *RawFile) ([]byte, error) {
			return data, nil
		},
	}
}

// NewRawTpl applies a go text template to the given tpl. The template is evaluated at render time.
// The use declaration is a no-op and just returns itself.
func NewRawTpl(name, mimeType string, n *Tpl) *RawFile {
	return &RawFile{
		Name:     name,
		MimeType: mimeType,
		Data: func(f *RawFile) ([]byte, error) {
			return f.renderText(n)
		},
	}
}

func (f *RawFile) renderText(n *Tpl) ([]byte, error) {
	ctx := &tplTextRenderContext{
		tpl: n,
	}

	tmpl, err := template.New(n.ObjPos.String()).Parse(n.Template)
	if err != nil {
		return nil, fmt.Errorf("cannot parse template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, ctx); err != nil {
		return nil, fmt.Errorf("cannot execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// ensure that we always implement the full contract
var _ TplContext = (*tplTextRenderContext)(nil)

type tplTextRenderContext struct {
	tpl *Tpl
}

func (t *tplTextRenderContext) Get(key string) interface{} {
	return t.tpl.Values[key]
}

func (t *tplTextRenderContext) Use(name string) string {
	return name
}

func (t *tplTextRenderContext) Self() *Tpl {
	return t.tpl
}
