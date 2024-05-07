package wdl

type TypeRef struct {
	Qualifier PkgImportQualifier
	Name      Identifier
	Pointer   bool
	Params    []*TypeRef
	TypeParam bool // TODO how to model bound (go) and upper and lower bounds (java)

}

func (t TypeRef) String() string {
	tmp := string(t.Qualifier) + "." + string(t.Name)
	if t.Pointer {
		tmp = "*" + tmp
	}
	if len(t.Params) > 0 {
		tmp += "["
		for _, param := range t.Params {
			tmp += param.String() + ", "
		}
		tmp += "]"
	}

	return tmp
}
