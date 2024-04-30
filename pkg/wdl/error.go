package wdl

import "fmt"

type ErrorWithPos struct {
	Pos   Pos
	Cause error
}

func NewErrorWithPos(pos Pos, cause error) *ErrorWithPos {
	return &ErrorWithPos{Pos: pos, Cause: cause}
}

func (e ErrorWithPos) Unwrap() error { return e.Cause }

func (e ErrorWithPos) Error() string {
	return fmt.Sprintf("pos:%v cause:%v", e.Pos, e.Cause)
}
