package wdl

// Identifier follows the typical identifier rules of Go or Java. Examples:
//   - Go: string, PkgName, Error, String, int64
//   - Java: List, String, Integer, int, toString
type Identifier string

func (id Identifier) String() string {
	return string(id)
}

func (id Identifier) expression() {}
