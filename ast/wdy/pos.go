package wdy

import "fmt"

type Pos struct {
	File string
	Line int
}

func (p Pos) String() string {
	return fmt.Sprintf("%s:%d", p.File, p.Line)
}
