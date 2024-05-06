package typescript

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
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
	if rtype.TypeParam() {
		return rtype
	}

	if rtype.Pkg() == nil {
		panic(fmt.Errorf("type %#v has empty package", rtype))
	}
	if rtype.Pkg().Qualifier() == r.selfImportPath || rtype.Pkg().Qualifier() == "std" {
		return rtype
	}

	r.namedImports[rtype.Name()] = rtype.Pkg().Qualifier() // todo this is wrong for collisions on name but different paths
	return rtype
}

func (r *RFile) AddImport(pkgName wdl.Identifier, path wdl.PkgImportQualifier) {
	r.namedImports[pkgName] = path
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

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	rFile := NewRFile(r, file.Pkg().Qualifier())
	tmp := &render.Writer{}
	for _, node := range file.TypeDefs() {
		if err := rFile.renderTypeDef(node, tmp); err != nil {
			return nil, err
		}
	}

	for name, qual := range file.Imports() {
		rFile.AddImport(name, qual)
	}

	if len(rFile.namedImports) > 0 {
		tmpImports := make(map[wdl.Identifier]wdl.PkgImportQualifier)
		for identifier, qualifier := range rFile.namedImports {
			tmpImports[identifier] = qualifier
		}
		for identifier, qualifier := range file.Imports() {
			tmpImports[identifier] = qualifier
		}

		for namedImport, qualifier := range tmpImports {
			w.Printf("import type { %s } from '%s';\n", tsUpperNameStr(string(namedImport)), atPathName(string(qualifier)))
		}

	}

	w.Printf("\n")

	w.Printf(tmp.String())

	return w.Bytes(), nil // TODO do we have something like autoformat in typescript?
}

func atPathName(p string) string {
	return p
}
