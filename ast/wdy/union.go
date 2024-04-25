package wdy

import "encoding/json"

type Union struct {
	Ref     TypeReference
	Macros  []Macro
	Comment []string
	Types   []TypeReference
	typeDecl
}

func (u *Union) GetMacros() []Macro {
	return u.Macros
}

func (u *Union) GetRef() TypeReference {
	return u.Ref
}

func (u *Union) String() string {
	buf, err := json.MarshalIndent(u, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
