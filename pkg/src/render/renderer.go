package render

import "github.com/worldiety/macro/pkg/src/ast"

// Artifact is either a Dir or File.
type Artifact interface {
	Name() string
}

// Renderer maps the given node into an artifact.
type Renderer interface {
	Render(node ast.Node) (Artifact, error)
}
