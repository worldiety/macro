package wdy

import "encoding/json"

type Interface struct {
	Ref     TypeReference
	Macros  []Macro
	Comment []string
	Methods []Func
	typeDecl
}

type Func struct {
	Name string
}

func (u *Interface) String() string {
	buf, err := json.MarshalIndent(u, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
