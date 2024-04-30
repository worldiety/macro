package wdl

import "fmt"

type Pos struct {
	File string
	Line int
}

func NewPos(file string, line int) Pos {
	return Pos{
		File: file,
		Line: line,
	}
}

func (p Pos) String() string {
	return fmt.Sprintf("%s:%d", p.File, p.Line)
}
