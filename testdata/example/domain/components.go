package domain

// A Component is a sum type or tagged union.
// Actually, we can generate different flavors, so that Go makes fun for modelling business stuff.
//
// #[go.TaggedUnion "json":"adjacent", "tag":"type", "names":["helloBtn","ATef","str","Chappy"]]
type _Component interface {
	Button | TextField | Text | Chapter
}

// #[go.TaggedUnion]
type _privateComponent interface {
	Button | TextField | Text | Chapter
}

// #[go.TaggedUnion]
type _ÜmlautComponent interface {
	Button | TextField | Text | Chapter | Äpfel
}

// #[go.TaggedUnion "json":"internal"]
type _ExampleType interface {
	Button | TextField
}

type Äpfel string

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
