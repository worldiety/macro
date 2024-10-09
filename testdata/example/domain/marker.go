package domain

// #[go.TaggedUnion "markerMethod":true]
type _Frucht interface {
	Birne | Kirsche
}

type Birne string
type Kirsche string

func blub() {
	var frucht Frucht
	frucht = NewFrucht(Birne("asd"))
	frucht = NewFrucht(Kirsche("x"))
	_ = frucht

}
