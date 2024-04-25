package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
	"strings"
)

func (r *Renderer) renderMod(mod *ast.Mod, parent *render.Dir) (*render.Dir, error) {
	modDir := r.ensurePkgDir(mod.Target.Out, parent)
	modDir.MimeType = MimeTypeGoModule

	modDir.Files = append(modDir.Files, &render.File{
		FileName: PackageGoModFile,
		MimeType: MimeTypeGoMod,
		Buf:      []byte(createGoModFile(mod)),
	})

	var firstErr error

	for _, pkg := range mod.Pkgs {
		// we cannot use name here, because in go the name and the import path may be different
		if !strings.HasPrefix(pkg.Path, mod.Name) {
			return nil, fmt.Errorf("declared package '%s' must be prefixed by module path '%s'", pkg.Path, mod.Name)
		}
		var pkgDir *render.Dir
		if pkg.Path == mod.Name {
			pkgDir = modDir
		} else {
			pkgDir = r.ensurePkgDir(pkg.Path[len(mod.Name)+1:], modDir)
		}

		files, err := r.renderPkg(pkg)
		if firstErr == nil && err != nil {
			firstErr = fmt.Errorf("cannot render package '%s': %w", pkg.Path, err)
		}

		pkgDir.Files = append(pkgDir.Files, files...)
	}

	return modDir, firstErr
}

func createGoModFile(mod *ast.Mod) string {
	var tmp strings.Builder

	tmp.WriteString(fmt.Sprintf("module %s\n\ngo %s\n", mod.Name, mod.Target.MinLangVersion))

	if len(mod.Target.Require.GoMod) > 0 {
		tmp.WriteString("\nrequire (\n")
		for _, dep := range mod.Target.Require.GoMod {
			tmp.WriteString("\t")
			tmp.WriteString(dep)
			tmp.WriteString("\n")
		}

		tmp.WriteString(")\n")
	}

	return tmp.String()
}
