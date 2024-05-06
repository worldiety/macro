package wdl

type TypeRef struct {
	Qualifier PkgImportQualifier
	Name      Identifier
	Pointer   bool
	Params    []*TypeRef
	TypeParam bool // TODO how to model bound (go) and upper and lower bounds (java)

}

func (t TypeRef) String() string {
	return string(t.Qualifier) + "." + string(t.Name)
}
