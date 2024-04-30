package golang

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"strconv"
)

type RFile struct {
	parent         *Renderer
	selfImportPath wdl.PkgImportQualifier
	namedImports   map[wdl.Identifier]wdl.PkgImportQualifier
}

func NewRFile(parent *Renderer, pkg wdl.PkgImportQualifier) *RFile {
	return &RFile{parent: parent, namedImports: map[wdl.Identifier]wdl.PkgImportQualifier{}, selfImportPath: pkg}
}

func (r *RFile) Use(rtype *wdl.ResolvedType) *wdl.ResolvedType {
	if rtype.Pkg().Qualifier() == r.selfImportPath || rtype.Pkg().Qualifier() == "std" {
		return rtype
	}

	r.namedImports[rtype.Pkg().Name()] = rtype.Pkg().Qualifier() // todo this is wrong for collisions on name but different paths
	return rtype
}

// RenderFile generates the code for the entire file.
func (r *Renderer) RenderFile(file *wdl.File) ([]byte, error) {
	w := &render.Writer{}

	// file license or whatever
	if file.Preamble() != nil {
		r.writeComment(w, false, file.Pkg().Name().String(), file.Preamble().String())
		w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	}

	// actual package comment
	if file.Comment() != nil {
		r.writeComment(w, true, file.Pkg().Name().String(), file.Comment().String())
	}

	w.Printf("package %s\n", file.Pkg().Name())

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	rFile := NewRFile(r, file.Pkg().Qualifier())
	tmp := &render.Writer{}
	for _, node := range file.TypeDefs() {
		if err := rFile.renderTypeDef(node, tmp); err != nil {
			return nil, err
		}
	}

	if len(rFile.namedImports) > 0 {
		w.Printf("import (\n")
		for namedImport, qualifier := range rFile.namedImports {
			w.Printf("  %s %s\n", namedImport, strconv.Quote(string(qualifier)))
		}

		w.Printf(")\n")
	}

	w.Printf(tmp.String())

	return Format(w.Bytes())
}
