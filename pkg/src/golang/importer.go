package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
)

// An importerKey declares a new type to secure private map access.
type importerKey int32

// lastImporterKey is a global thread-safe counter.
var lastImporterKey int32

func nextImporterKey() importerKey {
	return importerKey(atomic.AddInt32(&lastImporterKey, 1))
}

// importer manages the rendered import section at the files top.
type importer struct {
	selfImportPath string
	namedImports   map[string]string // named import => qualifier
}

// newImporter allocates an according instance.
func newImporter(selfImportPath string) *importer {
	return &importer{
		selfImportPath: selfImportPath,
		namedImports:   map[string]string{},
	}
}

// installImporter installs a new importer instance into every ast.SrcFileNode.
func installImporter(r *Renderer) error {
	r.importerId = nextImporterKey()
	return ast.ForEachMod(r.root, func(mod *ast.Mod) error {
		for _, pkg := range mod.Pkgs {
			for _, file := range pkg.PkgFiles {
				file.PutValue(r.importerId, newImporter(pkg.Path))
			}
		}

		return nil
	})
}

// installImporter overwrites all registered importers with a nil value.
func uninstallImporter(r *Renderer) error {
	return ast.ForEachMod(r.root, func(mod *ast.Mod) error {
		mod.PutValue(r.importerId, nil)

		return nil
	})
}

// importerFromTree walks up the tree until it finds the first importer from any ast.Node.Value.
func importerFromTree(r *Renderer, n ast.Node) *importer {
	root := n
	for root != nil {
		if imp, ok := root.Value(r.importerId).(*importer); ok {
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
	for _, name := range p.namedImports {
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
func (p *importer) shortify(name ast.Name) ast.Name {
	qual := name.Qualifier()
	id := name.Identifier()
	if id == "" || qual == "" {
		return name
	}

	if qual == p.selfImportPath {
		return ast.Name(id)
	}

	namedImportIdxName := strings.LastIndex(qual, "/") // e.g. 3 for net/http or -1 for net
	if namedImportIdxName == -1 {
		namedImportIdxName = 0
	}

	namedImport := MakePrivate(MakeIdentifier(qual[namedImportIdxName:]))

	otherQualifier, inScope := p.namedImports[namedImport]
	if inScope {
		// already registered the identical qualifier, e.g.
		// net/http => http
		if otherQualifier == qual {
			return ast.Name(namedImport + "." + id)
		} else {
			// name collision, build something artificial with increasing number
			num := 1
			for {
				num++
				namedImport2 := namedImport + strconv.Itoa(num)
				otherQualifier, inScope = p.namedImports[namedImport2]
				if inScope {
					if otherQualifier == qual {
						return ast.Name(namedImport + "." + id)
					}
					// loop again until either found or no other entry found
				} else {
					namedImport = namedImport2
					break
				}
			}
		}
	}

	p.namedImports[namedImport] = qual
	return ast.Name(namedImport + "." + id)
}
