package markdown

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"github.com/worldiety/macro/pkg/wdl/render"
	"path/filepath"
	"slices"
	"strings"
)

func (m *Markdown) makeDoc(opts markdownParams, def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) *render.Writer {
	md := &render.Writer{}
	m.chapterProject(md)

	boundedContexts := m.collectMergedBoundedContexts(m.prog.Annotations())
	for _, bctx := range boundedContexts {
		md.Printf("## %s\n\n", bctx.name)
		text := ""
		for _, bc := range bctx.bcs {
			text += commentText1(bc.Pkg()) + "\n"
		}
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			md.Print("Dieser Bounded Context ist noch nicht dokumentiert.\n\n")
		} else {
			md.Print(text)
		}
		md.Print("\n\n")

		if len(bctx.useCases) > 0 {
			md.Printf("### Anwendungsfälle\n\n")
			for _, usecase := range bctx.useCases {
				md.Printf("#### %s\n\n", usecase.Name())
				text := commentText1(usecase.Fn())
				if text == "" {
					md.Print("Dieser Anwendungsfall ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}
		}

		if len(bctx.useCases) > 0 {
			md.Printf("### Werte\n\n")
			for _, a := range bctx.values {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieser Werttyp ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}
		}

		if len(bctx.entities) > 0 {
			md.Printf("### Entitäten\n\n")
			for _, a := range bctx.entities {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Diese Entität ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}
		}

		if len(bctx.aggregates) > 0 {
			md.Printf("### Aggregate\n\n")
			for _, a := range bctx.aggregates {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieses Aggregat ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}
		}
		if len(bctx.events) > 0 {
			md.Printf("### Domänenereignisse\n\n")
			for _, a := range bctx.events {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieses Ereignis ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}
		}

		if len(bctx.services) > 0 {
			md.Printf("### Domänenservices\n\n")
			for _, a := range bctx.services {
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

		if len(bctx.repos) > 0 {
			md.Printf("### Repositories\n\n")
			for _, a := range bctx.repos {
				md.Printf("#### %s\n\n", a.Name())
				text := commentText1(a.TypeDef())
				if text == "" {
					md.Print("Dieses Repository ist noch nicht dokumentiert.\n\n")
				} else {
					md.Print(text)
					md.Print("\n\n")
				}
			}
		}
	}

	if !opts.OmitSecurityChapter {
		m.chapterSecurity(md)
	}

	return md
}

func (m *Markdown) chapterProject(md *render.Writer) {
	var prj *wdl.ProjectAnnotation
	for _, annotation := range m.prog.Annotations() {
		if a, ok := annotation.(*wdl.ProjectAnnotation); ok {
			prj = a
			break
		}
	}

	prjName := filepath.Base(m.prog.Path())
	if prj != nil {
		prjName = prj.Name()
	}

	md.Printf("# %s\n\n", prjName)
	if prj == nil || prj.Pkg().Comment().String() == "" {
		md.Printf("Das Projekt ist noch nicht dokumentiert.\n\n")
	} else {
		md.Print(commentText1(prj.Pkg()))
		md.Print("\n\n")
	}
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
				md.Printf("|%s|%s|\n", permission.PermissionID(), linkify(m.alias(permission.TypeDef()), m.alias(permission.TypeDef())))
			}
		}

		// there is a later post macro which validates the actual uniqueness
		permissions = slices.CompactFunc(permissions, func(annotation *wdl.PermissionAnnotation, annotation2 *wdl.PermissionAnnotation) bool {
			return annotation.PermissionID() == annotation2.PermissionID()
		})

		md.Printf("\n")

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
				md.Printf("|jeder|%s|\n", linkify(uc.Name(), uc.Name()))
			}
		}

	} else {
		md.Printf("Es sind noch keine Anwendungsfälle definiert.\n\n")
	}
}

func linkify(name, target string) string {
	target = strings.ToLower(target)
	target = strings.ReplaceAll(target, " ", "-")
	return fmt.Sprintf("[%s](#%s)", name, target)
}

type aggregatedBC struct {
	name       string
	bcs        []*wdl.BoundedContextAnnotation
	useCases   []*wdl.UseCaseAnnotation
	repos      []*wdl.RepositoryAnnotation
	values     []*wdl.ValueAnnotation
	entities   []*wdl.EntityAnnotation
	aggregates []*wdl.AggregateRootAnnotation
	events     []*wdl.DomainEventAnnotation
	services   []*wdl.DomainServiceAnnotation
}

func (m *Markdown) collectMergedBoundedContexts(annotations []wdl.Annotation) []aggregatedBC {
	tmp := make(map[string]aggregatedBC)
	for _, annotation := range annotations {
		if bcAn, ok := annotation.(*wdl.BoundedContextAnnotation); ok {
			abc := tmp[bcAn.Name()]
			abc.name = bcAn.Name()
			abc.bcs = append(abc.bcs, bcAn)
			slices.SortFunc(abc.bcs, func(a, b *wdl.BoundedContextAnnotation) int {
				// provide a stable order
				return strings.Compare(string(a.Pkg().Qualifier()), string(b.Pkg().Qualifier()))
			})

			// use cases
			for _, usecase := range m.collectUsecases(bcAn.Pkg()) {
				abc.useCases = append(abc.useCases, usecase)
			}

			slices.SortFunc(abc.useCases, func(a, b *wdl.UseCaseAnnotation) int {
				return strings.Compare(a.Name(), b.Name())
			})

			abc.useCases = slices.Compact(abc.useCases)

			// repos
			for _, repo := range collectWithTypeDef[*wdl.RepositoryAnnotation](m, bcAn.Pkg()) {
				abc.repos = append(abc.repos, repo)
			}

			slices.SortFunc(abc.repos, func(a, b *wdl.RepositoryAnnotation) int {
				return strings.Compare(a.Name(), b.Name())
			})

			abc.repos = slices.Compact(abc.repos)

			// values
			for _, value := range collectWithTypeDef[*wdl.ValueAnnotation](m, bcAn.Pkg()) {
				abc.values = append(abc.values, value)
			}

			slices.SortFunc(abc.values, func(a, b *wdl.ValueAnnotation) int {
				return strings.Compare(a.Name(), b.Name())
			})

			abc.values = slices.Compact(abc.values)

			// entities
			for _, ent := range collectWithTypeDef[*wdl.EntityAnnotation](m, bcAn.Pkg()) {
				abc.entities = append(abc.entities, ent)
			}

			slices.SortFunc(abc.entities, func(a, b *wdl.EntityAnnotation) int {
				return strings.Compare(a.Name(), b.Name())
			})

			abc.entities = slices.Compact(abc.entities)

			// aggregates
			for _, agr := range collectWithTypeDef[*wdl.AggregateRootAnnotation](m, bcAn.Pkg()) {
				abc.aggregates = append(abc.aggregates, agr)
			}

			slices.SortFunc(abc.aggregates, func(a, b *wdl.AggregateRootAnnotation) int {
				return strings.Compare(a.Name(), b.Name())
			})

			abc.aggregates = slices.Compact(abc.aggregates)

			// events
			for _, evt := range collectWithTypeDef[*wdl.DomainEventAnnotation](m, bcAn.Pkg()) {
				abc.events = append(abc.events, evt)
			}

			slices.SortFunc(abc.events, func(a, b *wdl.DomainEventAnnotation) int {
				return strings.Compare(a.Name(), b.Name())
			})

			abc.events = slices.Compact(abc.events)

			// services
			for _, srv := range collectWithTypeDef[*wdl.DomainServiceAnnotation](m, bcAn.Pkg()) {
				abc.services = append(abc.services, srv)
			}

			slices.SortFunc(abc.services, func(a, b *wdl.DomainServiceAnnotation) int {
				return strings.Compare(a.Name(), b.Name())
			})

			abc.services = slices.Compact(abc.services)

			// update slot
			tmp[bcAn.Name()] = abc
		}
	}

	var result []aggregatedBC
	for _, bc := range tmp {
		result = append(result, bc)
	}

	slices.SortFunc(result, func(a, b aggregatedBC) int {
		return strings.Compare(a.name, b.name)
	})

	return result
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

	slices.SortFunc(tmp, func(a, b *wdl.UseCaseAnnotation) int {
		return strings.Compare(a.Name(), b.Name())
	})
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
