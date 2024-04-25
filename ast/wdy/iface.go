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

func (u *Interface) GetMacros() []Macro {
	return u.Macros
}

func (u *Interface) GetRef() TypeReference {
	return u.Ref
}

func (u *Interface) String() string {
	buf, err := json.MarshalIndent(u, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
