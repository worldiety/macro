package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/render"
)

const (
	PackageGoDocFile = "doc.go"
	PackageGoModFile = "go.mod"

	MimeTypeGo       = "text/x-go-source"
	MimeTypeGoMod    = "text/x-go-source-mod"
	MimeTypeDir      = "application/x-directory"
	MimeTypeGoModule = "application/x-directory-module"
)

// Options for the renderer.
type Options struct {
}

// Renderer provides a go renderer.
type Renderer struct {
	opts       Options
	root       ast.Node
	importerId importerKey // track our unique key, to perform cleanup
}

// NewRenderer creates a new Renderer instance.
func NewRenderer(opts Options) *Renderer {
	return &Renderer{opts: opts}
}

// tearUp prepares the ast to be used for source generation.
func (r *Renderer) tearUp(node ast.Node) error {
	r.root = ast.Root(node)

	if err := installImporter(r); err != nil {
		return fmt.Errorf("unable to install importer: %w", err)
	}

	return nil
}

// tearDown frees allocated resources.
func (r *Renderer) tearDown() error {
	if err := uninstallImporter(r); err != nil {
		return fmt.Errorf("unable to uninstall importer: %w", err)
	}

	return nil
}

// importer resolves the current importer from the parents file.
func (r *Renderer) importer(n ast.Node) *importer {
	return importerFromTree(r, n)
}

// Render converts the given node into a render.Artifact. A partial result is returned if an error is detected.
func (r *Renderer) Render(node ast.Node) (a render.Artifact, err error) {
	if err := r.tearUp(node); err != nil {
		return nil, fmt.Errorf("unable to tearUp: %w", err)
	}

	defer func() {
		if e := r.tearDown(); e != nil && err == nil {
			err = e
		}
	}()

	root := &render.Dir{}
	err = ast.ForEachMod(node, func(mod *ast.Mod) error {
		if mod.Target.Lang == ast.LangGo {
			_, err := r.renderMod(mod, root)

			if err != nil {
				return fmt.Errorf("cannot render module '%s': %w", mod.Name, err)
			}
		}

		return nil
	})

	if err != nil {
		return root, fmt.Errorf("cannot render project: %w", err)
	}

	return root, nil
}
