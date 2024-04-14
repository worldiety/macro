package wdy

import "encoding/json"

type Union struct {
	Ref     TypeReference
	Macros  []Macro
	Comment []string
	Types   []TypeReference
	typeDecl
}

func (u *Union) String() string {
	buf, err := json.MarshalIndent(u, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
