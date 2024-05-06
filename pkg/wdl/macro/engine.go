package macro

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render/golang"
	"github.com/worldiety/macro/pkg/wdl/render/typescript"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Engine struct {
	prog                                *wdl.Program
	macros                              map[wdl.MacroName]func(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error
	preamble                            string
	aggregateGeneratedGoFilesPerPackage bool
}

func NewEngine(pg *wdl.Program) *Engine {
	e := &Engine{prog: pg, preamble: "Code generated by github.com/worldiety/macro. DO NOT EDIT.", aggregateGeneratedGoFilesPerPackage: true}
	e.macros = map[wdl.MacroName]func(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error{
		"go.TaggedUnion": e.goTaggedUnion,
		"go.TypeScript":  e.goIntoTypescript,
	}
	return e
}

func (e *Engine) Exec() error {
	for _, pkg := range e.prog.Packages() {
		for _, def := range pkg.TypeDefs() {
			for _, macroInvoc := range def.Macros() {
				macro := e.macros[macroInvoc.Name()]
				if macro == nil {
					return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("macro %s not found", macroInvoc.Name()))
				}

				if err := macro(def, macroInvoc); err != nil {
					return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("macro %s execution error: %w", macroInvoc.Name(), err))
				}
			}
		}
	}

	return nil
}

// Emit writes all changed files into their absolute path positions.
func (e *Engine) Emit() error {
	if e.aggregateGeneratedGoFilesPerPackage {
		for _, pkg := range e.prog.Packages() {
			var tmp []*wdl.File
			aggregatedFile := wdl.NewFile(nil)
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
			}

		}
	}

	return nil
}

func stripName(s string) string {
	return strings.TrimLeft(s, "._")
}
