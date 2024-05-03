package domain

import "example/domain/xcompo"

// A Component is a sum type or tagged union.
// Actually, we can generate different flavors, so that Go makes fun for modelling business stuff.
//
// #[go.TaggedUnion "json":"intern", "tag":"type"]
type _Component interface {
	Button | TextField | Text | Chapter | xcompo.RichText | xcompo.Icon | string | []string | []Text | map[int]Button
}

type Blub error

// Another doc.
type Another interface {
	// JustForDoc is stuff.
	JustForDoc()
	error
}

type Chapter int

type Button struct {
	Caption string
	T       error
}
