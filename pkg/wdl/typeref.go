package wdl

type TypeRef struct {
	Qualifier PkgImportQualifier
	Name      Identifier
}

func (t TypeRef) String() string {
	return string(t.Qualifier) + "." + string(t.Name)
}
