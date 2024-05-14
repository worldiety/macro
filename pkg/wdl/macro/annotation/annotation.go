package annotation

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
)

const (
	entityAnnotation         = "@Entity"
	aggregateAnnotation      = "@Aggregate"
	valueAnnotation          = "@Value"
	domainEventAnnotation    = "@DomainEvent"
	usecaseAnnotation        = "@Usecase"
	repositoryAnnotation     = "@Repository"
	boundedContextAnnotation = "@BoundedContext"
	domainService            = "@DomainService"
	project                  = "@Project"
)

type Annotation struct {
	prog     *wdl.Program
	preamble string
}

func NewAnnotation(prog *wdl.Program, preamble string) *Annotation {
	return &Annotation{prog: prog, preamble: preamble}
}

func (m *Annotation) Names() []wdl.MacroName {
	return []wdl.MacroName{
		entityAnnotation,
		aggregateAnnotation,
		valueAnnotation,
		domainEventAnnotation,
		usecaseAnnotation,
		boundedContextAnnotation,
		repositoryAnnotation,
		domainService,
		project,
	}
}

func (m *Annotation) Expand(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {
	var a wdl.Annotation
	switch macroInvoc.Name() {
	case entityAnnotation:
		a = wdl.NewEntityAnnotation(macroInvoc.Value(), macroInvoc.Pos(), def)
	case aggregateAnnotation:
		a = wdl.NewAggregateRootAnnotation(macroInvoc.Value(), macroInvoc.Pos(), def)
	case valueAnnotation:
		a = wdl.NewValueAnnotation(macroInvoc.Value(), macroInvoc.Pos(), def)
	case domainEventAnnotation:
		a = wdl.NewValueAnnotation(macroInvoc.Value(), macroInvoc.Pos(), def)
	case repositoryAnnotation:
		a = wdl.NewRepositoryAnnotation(macroInvoc.Value(), macroInvoc.Pos(), def)
	case domainService:
		a = wdl.NewDomainServiceAnnotation(macroInvoc.Value(), macroInvoc.Pos(), def)
	case usecaseAnnotation:
		if fn, ok := def.(*wdl.Func); ok {
			a = wdl.NewUseCaseAnnotation(macroInvoc.Value(), fn, macroInvoc.Pos())
		} else {
			return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("only functions or methods can be annotated as usecase"))
		}
	case boundedContextAnnotation:
		if pkg, ok := def.(*wdl.Package); ok {
			a = wdl.NewBoundedContextAnnotation(macroInvoc.Value(), pkg, macroInvoc.Pos())
		} else {
			return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("only packages can be annotated as bounded context"))
		}
	case project:
		if pkg, ok := def.(*wdl.Package); ok {
			a = wdl.NewProjectAnnotation(macroInvoc.Value(), pkg, macroInvoc.Pos())
		} else {
			return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("only packages can be annotated as a project"))
		}
	default:
		return wdl.NewErrorWithPos(macroInvoc.Pos(), fmt.Errorf("unknown annotation type"))
	}

	m.prog.AddAnnotations(a)
	return nil
}
