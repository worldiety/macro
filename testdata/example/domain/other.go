package domain

type Text string

type TextField struct {
	Label    Text
	Hint     string
	MaxChars int32
}
