package macro

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render/typescript"
	"path/filepath"
	"regexp"
)

type goIntoTypescriptParams struct {
	Path string `json:"path"`
}

var regexSrc = regexp.MustCompile(`.*src/`)

type TranspileTypeScript struct {
	prog     *wdl.Program
	preamble string
}

func NewTranspileTypeScript(prog *wdl.Program, preamble string) *TranspileTypeScript {
	return &TranspileTypeScript{prog: prog, preamble: preamble}
}

func (m *TranspileTypeScript) Names() []wdl.MacroName {
	return []wdl.MacroName{"go.TypeScript", "TypeScript"}
}

func (m *TranspileTypeScript) Expand(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {

	var params goIntoTypescriptParams
	if err := macroInvoc.UnmarshalParams(&params); err != nil {
		return err
	}

	tsDir := regexSrc.ReplaceAllString(params.Path, "@/")

	tspkg, ok := m.prog.PackageByPath(wdl.PkgImportQualifier(tsDir))
	if !ok {
		tspkg = wdl.NewPackage(func(pkg *wdl.Package) {
			pkg.SetQualifier(wdl.PkgImportQualifier(tsDir))
		})
		m.prog.AddPackage(tspkg)
	}

	tspkg.AddFiles(wdl.NewFile(func(file *wdl.File) {
		file.SetMimeType(wdl.MimeTypeTypeScript)
		file.SetPath(filepath.Join(m.prog.Path(), params.Path))
		file.SetName(typescript.GetFilename(def.Name()))
		file.SetModified(true)
		file.SetGenerated(true)
		file.SetPreamble(wdl.NewComment(func(comment *wdl.Comment) {
			comment.AddLines(wdl.NewCommentLine(func(line *wdl.CommentLine) {
				line.SetText(m.preamble)
			}))
		}))

		file.AddTypeDefs(def)
	}))

	return nil
}
