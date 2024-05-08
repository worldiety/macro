package macro

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/macro/annotation"
	"github.com/worldiety/macro/pkg/wdl/macro/markdown"
	"github.com/worldiety/macro/pkg/wdl/render/golang"
	"github.com/worldiety/macro/pkg/wdl/render/typescript"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Macro interface {
	Names() []wdl.MacroName
	Expand(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error
}

type Engine struct {
	prog                                *wdl.Program
	macros                              []Macro
	preamble                            string
	aggregateGeneratedGoFilesPerPackage bool
}

func NewEngine(pg *wdl.Program) *Engine {
	e := &Engine{prog: pg, preamble: "Code generated by github.com/worldiety/macro. DO NOT EDIT.", aggregateGeneratedGoFilesPerPackage: true}
	e.macros = []Macro{
		annotation.NewAnnotation(pg, e.preamble),
		NewGoTaggedUnion(pg, e.preamble),
		NewTranspileTypeScript(pg, e.preamble),
		markdown.NewMarkdown(pg, e.preamble),
	}

	return e
}

func (e *Engine) Exec() error {
	unusedMacroDeclarations := map[*wdl.MacroInvocation]struct{}{}
	for _, pkg := range e.prog.Packages() {
		for _, macroInvoc := range pkg.Comment().Macros() {
			unusedMacroDeclarations[macroInvoc] = struct{}{}
		}

		for _, def := range pkg.TypeDefs() {
			for _, macroInvoc := range def.Comment().Macros() {
				unusedMacroDeclarations[macroInvoc] = struct{}{}
			}
		}
	}

	// we execute macros in our build-in order, e.g. to apply annotations and code expanders before exporters.
	// independence of declaration order, makes our (intended) side effects more reasonable.
	for _, macro := range e.macros {
		for _, pkg := range e.prog.Packages() {
			for _, macroInvoc := range pkg.Comment().Macros() {
				for _, name := range macro.Names() {
					if name == macroInvoc.Name() {
						delete(unusedMacroDeclarations, macroInvoc)
						if err := macro.Expand(pkg, macroInvoc); err != nil {
							return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("macro %s execution error: %w", macroInvoc.Name(), err))
						}
					}
				}

			}
			for _, def := range pkg.TypeDefs() {
				for _, macroInvoc := range def.Comment().Macros() {
					for _, name := range macro.Names() {
						if name == macroInvoc.Name() {
							delete(unusedMacroDeclarations, macroInvoc)
							if err := macro.Expand(def, macroInvoc); err != nil {
								return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("macro %s execution error: %w", macroInvoc.Name(), err))
							}
						}
					}
				}
			}
		}
	}

	// check, if we were exhaustive
	var someUnused *wdl.MacroInvocation
	for unusedMacro := range unusedMacroDeclarations {
		someUnused = unusedMacro
		slog.Error("unknown macro", "name", unusedMacro.Name(), "pos", unusedMacro.Pos())
	}

	if len(unusedMacroDeclarations) > 0 {
		return wdl.NewErrorWithPos(someUnused.Pos(), fmt.Errorf("unknown macro"))
	}

	return nil
}

// Emit writes all changed files into their absolute path positions.
func (e *Engine) Emit() error {
	if e.aggregateGeneratedGoFilesPerPackage {
		for _, pkg := range e.prog.Packages() {
			var tmp []*wdl.File
			aggregatedFile := wdl.NewFile(nil)
			aggregatedFile.SetMimeType(wdl.MimeTypeGo)
			for _, file := range pkg.Files() {
				if file.Modified() && file.Generated() && file.MimeType() == wdl.MimeTypeGo {
					aggregatedFile.Import(file)
				} else {
					tmp = append(tmp, file)
				}
			}

			aggregatedFile.SetName("macros.gen.go")
			tmp = append(tmp, aggregatedFile)
			pkg.SetFiles(tmp)

		}
	}

	for _, pkg := range e.prog.Packages() {
		for _, file := range pkg.Files() {
			if !file.Modified() {
				continue
			}
			switch file.MimeType() {
			case wdl.MimeTypeGo:

				renderer := golang.NewRenderer(golang.Options{})
				buf, err := renderer.RenderFile(file)
				if err != nil {
					fmt.Println(string(buf))
					return err
				}
				fname := filepath.Join(file.Path(), stripName(file.Name()))
				if err := os.WriteFile(fname, buf, os.ModePerm); err != nil {
					return err
				}

				slog.Info("emit go source", "file", fname)
			case wdl.MimeTypeTypeScript:
				renderer := typescript.NewRenderer(typescript.Options{})
				buf, err := renderer.RenderFile(file)
				if err != nil {
					fmt.Println(string(buf))
					return err
				}
				fname := filepath.Join(file.Path(), stripName(file.Name()))
				if err := os.WriteFile(fname, buf, os.ModePerm); err != nil {
					return err
				}

				slog.Info("emit typescript source", "file", fname)
			case wdl.Raw:
				fname := filepath.Join(file.Path(), stripName(file.Name()))
				if err := os.WriteFile(fname, file.RawBytes(), os.ModePerm); err != nil {
					return err
				}
				slog.Info("emit raw file", "file", fname)
			default:
				slog.Error("unknown mimetype file", "file", file.Name(), "mimetype", file.MimeType())
			}

		}
	}

	return nil
}

func stripName(s string) string {
	return strings.TrimLeft(s, "._")
}
