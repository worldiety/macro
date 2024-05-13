package markdown

import (
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"path/filepath"
	"slices"
	"strings"
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
				md.Printf("#### <a name=\"%s\"></a> %s\n\n", usecase.Fn().Name(), usecase.Name())
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

	m.chapterSecurity(md)

	return md
}

func (m *Markdown) chapterSecurity(md *render.Writer) {
	md.Printf("## Berechtigungskonzept\n\n")
	useCases := collectUseCases(m.prog.Annotations())
	if len(useCases) > 0 {
		md.Print("Im Folgenden werden alle auditierten Berechtigungen dargestellt.\nDiese Berechtigungen sind Aktor-gebunden, d.h. ein Nutzer oder Drittsysteme müssen diese Rechte zugewiesen bekommen haben, um den Anwendungsfall ausführen zu dürfen.\nAchtung: es kann dynamische bzw. objektbezogene Rechte in Anwendungsfällen geben, die unabhängig von Berechtigungen das Darstellen, Löschen oder Bearbeiten von vertraulichen Informationen erlaubt. Diese sind hier nicht erfasst, sondern sind in der jeweiligen Dokumentation der Anwendungsfälle erwähnt.\n\n")
		permissions := collectPermissions(m.prog.Annotations())
		if len(permissions) > 0 {
			slices.SortFunc(permissions, func(a, b *wdl.PermissionAnnotation) int {
				return strings.Compare(a.PermissionID(), b.PermissionID())
			})

			md.Printf("|Berechtigung|Anwendungsfall|\n|----|----|\n")
			for _, permission := range permissions {
				md.Printf("|%s|[%s](#%s)|\n", permission.PermissionID(), m.alias(permission.TypeDef()), permission.TypeDef().Name())
			}
		}

		var anonUseCases []*wdl.UseCaseAnnotation
		for _, useCase := range useCases {
			found := false
			for _, permission := range permissions {
				if permission.TypeDef() == useCase.Fn() {
					found = true
					break
				}
			}

			if !found {
				anonUseCases = append(anonUseCases, useCase)
			}
		}

		if len(anonUseCases) > 0 {
			md.Printf("Die folgenden Anwendungsfälle sind grundsätzlich ohne Autorisierung verwendbar, erfordern also keine Berechtigungen und werden auch nicht auditiert.\n\n")
			md.Printf("|Berechtigung|Anwendungsfall|\n|----|----|\n")
			for _, uc := range anonUseCases {
				md.Printf("|jeder|[%s](#%s)|\n", m.alias(uc.Fn()), uc.Fn().Name())
			}
		}

	} else {
		md.Printf("Es sind noch keine Anwendungsfälle definiert.\n\n")
	}
}

func collectPermissions(annotations []wdl.Annotation) []*wdl.PermissionAnnotation {
	var tmp []*wdl.PermissionAnnotation
	for _, annotation := range annotations {
		if pAn, ok := annotation.(*wdl.PermissionAnnotation); ok {
			tmp = append(tmp, pAn)
		}
	}
	return tmp
}

func collectUseCases(annotations []wdl.Annotation) []*wdl.UseCaseAnnotation {
	var tmp []*wdl.UseCaseAnnotation
	for _, annotation := range annotations {
		if pAn, ok := annotation.(*wdl.UseCaseAnnotation); ok {
			tmp = append(tmp, pAn)
		}
	}
	return tmp
}

func collectWithTypeDef[T interface {
	TypeDef() wdl.TypeDef
	Name() string
}](m *Markdown, pkg *wdl.Package) []T {
	var res []T
	for _, annotation := range m.prog.Annotations() {
		if a, ok := annotation.(T); ok {
			if a.TypeDef().Pkg() == pkg {
				res = append(res, a)
			}
		}
	}

	slices.SortFunc(res, func(a, b T) int {
		return strings.Compare(a.Name(), b.Name())
	})
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

	slices.SortFunc(res, func(a, b *wdl.UseCaseAnnotation) int {
		return strings.Compare(a.Name(), b.Name())
	})

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

func (m *Markdown) alias(def wdl.TypeDef) string {
	for _, annotation := range m.prog.Annotations() {
		switch a := annotation.(type) {
		case *wdl.UseCaseAnnotation:
			if a.Fn() == def {
				return a.Name()
			}
		}
	}

	return def.Name().String()
}
