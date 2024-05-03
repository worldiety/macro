package domain

// A Component is a sum type or tagged union.
// Actually, we can generate different flavors, so that Go makes fun for modelling business stuff.
//
// #[go.TaggedUnion "json":"intern", "tag":"type", "names":["helloBtn","ATef","str","Chappy"]]
type _Component interface {
	Button | TextField | Text | Chapter
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
