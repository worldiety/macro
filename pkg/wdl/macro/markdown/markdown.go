package markdown

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"path/filepath"
)

type markdownParams struct {
	Out                 string `json:"out"` // the output markdown file. Default is <module root>/README.md
	OmitSecurityChapter bool   `json:"omitSecurityChapter"`
}

type Markdown struct {
	prog     *wdl.Program
	preamble string
}

func NewMarkdown(prog *wdl.Program, preamble string) *Markdown {
	return &Markdown{prog: prog, preamble: preamble}
}

func (m *Markdown) Names() []wdl.MacroName {
	return []wdl.MacroName{"markdown"}
}

func (m *Markdown) Expand(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {
	var opts markdownParams
	if err := macroInvoc.UnmarshalParams(&opts); err != nil {
		return fmt.Errorf("invalid macro params: %w", err)
	}

	if opts.Out == "" {
		opts.Out = "README.md"
	}

	opts.Out = filepath.Join(m.prog.Path(), opts.Out)

	md := m.makeDoc(opts, def, macroInvoc)

	m.prog.AddPackage(wdl.NewPackage(func(pkg *wdl.Package) {
		pkg.SetName("markdown")
		pkg.AddFiles(wdl.NewFile(func(file *wdl.File) {
			file.SetMimeType(wdl.Raw)
			file.SetName(filepath.Base(opts.Out))
			file.SetPath(filepath.Dir(opts.Out))
			file.SetRawBytes(md.Bytes())
			file.SetGenerated(true)
			file.SetModified(true)
		}))
	}))

	return nil
}
