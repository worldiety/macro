package java

import (
	"github.com/worldiety/macro/pkg/src"
	"github.com/worldiety/macro/pkg/src/ast"
	"sort"
)

type importerKey int

const importerId importerKey = 1

// importer manages the rendered import section at the files top.
type importer struct {
	identifiersInScope map[string]src.Name
}

func newImporter() *importer {
	return &importer{
		identifiersInScope: map[string]src.Name{},
	}
}

// installImporter installs a new importer instance into every ast.SrcFileNode.
func installImporter(n *ast.ModNode) {
	for _, node := range n.Packages() {
		for _, fileNode := range node.Files() {
			fileNode.SetValue(importerId, newImporter())
		}
	}
}

// importerFromTree walks up the tree until it finds the first importer from any ast.Node.Value.
func importerFromTree(n ast.Node) *importer {
	root := n
	for root != nil {
		if imp, ok := root.Value(importerId).(*importer); ok {
			return imp
		}

		newRoot := root.Parent()
		if newRoot == nil {
			panic("no attached importer found in ast scope")
		}

		root = newRoot
	}

	panic("invalid node")
}

// qualifiers returns the unique imported qualifiers.
func (p *importer) qualifiers() []string {
	tmp := map[string]string{}
	for _, name := range p.identifiersInScope {
		tmp[string(name)] = ""
	}

	var sorted []string
	for uniqueQualifier := range tmp {
		sorted = append(sorted, uniqueQualifier)
	}

	sort.Strings(sorted)

	return sorted
}

// shortify returns a qualified name, which is only valid in the importers scope. It may also decide to not import
// the given name, e.g. if a collision has been detected. If the name is a universe type or not complete, the original
// name is just returned.
func (p *importer) shortify(name src.Name) src.Name {
	qual := name.Qualifier()
	id := name.Identifier()
	if id == "" || qual == "" {
		return name
	}

	otherName, inScope := p.identifiersInScope[id]
	if inScope {
		// already registered the identical qualifier, e.g.
		// a.A => A
		// a.B => B
		if otherName == name {
			return src.Name(id)
		} else {
			// name collision
			return name
		}
	}

	p.identifiersInScope[id] = name
	return src.Name(id)
}
