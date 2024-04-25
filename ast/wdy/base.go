package wdy

import "path/filepath"

type TypeReference struct {
	Path     string          // Path to the package where this Type resides. Empty, if its a Basic type or from the default scope.
	Name     string          // Name of the actual type.
	TypeArgs []TypeReference // Type Parameter aka Generics
	typeDecl
}

func (t TypeReference) PathName() string {
	// TODO this is wrong in Go, because the import path has no relation to the actual package name
	return filepath.Base(t.Path)
}

func (t TypeReference) String() string {
	return t.Path + "." + t.Name
}

type typeDecl interface {
	isTypeDecl()
	GetRef() TypeReference
	GetMacros() []Macro
}

type TypeDecl typeDecl
