package domain

type Text string

type TextField struct {
	Label    Text
	Hint     string
	MaxChars int32
}

type Seq[V any] func(yield func(V) bool)

func Test(it Seq[string]) Seq[string] {
	return nil
}
