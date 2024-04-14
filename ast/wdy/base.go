package wdy

type TypeReference struct {
	Path     string          // Path to the package where this Type resides. Empty, if its a Basic type or from the default scope.
	Name     string          // Name of the actual type.
	TypeArgs []TypeReference // Type Parameter aka Generics
	typeDecl
}

func (t TypeReference) String() string {
	return t.Path + "." + t.Name
}

type typeDecl interface {
	isTypeDecl()
}

type TypeDecl typeDecl
