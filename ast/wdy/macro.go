package wdy

// Macro is text/template invocation like !{{go.TaggedUnion `"tag":"$_type"`}}.
// However, it is not used to generate text by itself.
// Instead, it contains only the execution commands and we (mis-)use the templating engine implementation for it.
type Macro struct {
	Template string
	Origin   Pos
}
