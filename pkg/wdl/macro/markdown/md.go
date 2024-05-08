package markdown

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"path/filepath"
)

func (m *Markdown) makeDoc(opts markdownParams, def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) *render.Writer {
	md := &render.Writer{}
	md.Printf("# %s\n\n", filepath.Base(m.prog.Path()))
	// TODO we need a project name and project preamble text to include here
	for _, annotation := range m.prog.Annotations() {
		if bc, ok := annotation.(*wdl.BoundedContextAnnotation); ok {
			md.Printf("## %s\n\n", bc.Name())

			md.Printf("### Anwendungsfälle\n\n")
			for _, usecase := range m.collectUsecases(bc.Pkg()) {
				md.Printf("#### %s\n\n", usecase.Name())
				text := commentText1(usecase.Fn())
				if text == "" {
					md.Print("Dieser Anwendungsfall ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}

			md.Printf("### Werte\n\n")
			for _, a := range collectWithTypeDef[*wdl.ValueAnnotation](m, bc.Pkg()) {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieser Werttyp ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}

			md.Printf("### Entitäten\n\n")
			for _, a := range collectWithTypeDef[*wdl.EntityAnnotation](m, bc.Pkg()) {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Diese Entität ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}

			md.Printf("### Aggregate\n\n")
			for _, a := range collectWithTypeDef[*wdl.AggregateRootAnnotation](m, bc.Pkg()) {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieses Aggregat ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}

			md.Printf("### Domänenereignisse\n\n")
			for _, a := range collectWithTypeDef[*wdl.DomainEventAnnotation](m, bc.Pkg()) {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieses Ereignis ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}

			md.Printf("### Domänenservices\n\n")
			for _, a := range collectWithTypeDef[*wdl.DomainServiceAnnotation](m, bc.Pkg()) {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieser Service ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}

		}
	}

	return md
}

func collectWithTypeDef[T interface{ TypeDef() wdl.TypeDef }](m *Markdown, pkg *wdl.Package) []T {
	var res []T
	for _, annotation := range m.prog.Annotations() {
		if a, ok := annotation.(T); ok {
			if a.TypeDef().Pkg() == pkg {
				res = append(res, a)
			}
		}
	}
	return res
}

func (m *Markdown) collectUsecases(pkg *wdl.Package) []*wdl.UseCaseAnnotation {
	var res []*wdl.UseCaseAnnotation
	for _, annotation := range m.prog.Annotations() {
		if a, ok := annotation.(*wdl.UseCaseAnnotation); ok {
			if a.Fn().Pkg() == pkg {
				res = append(res, a)
			}
		}
	}
	return res
}

func commentText1(e interface{ Comment() *wdl.Comment }) string {
	if e.Comment() == nil {
		return ""
	}

	tmp := ""
	for _, line := range e.Comment().Lines() {
		tmp += line.Text() + "\n"
	}
	return tmp
}
