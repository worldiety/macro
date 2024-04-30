package wdl

// File represents a physical source code file respective compilation unit.
//   - Go: <lowercase AnnotationName>.go
//   - Java: <CamelCasePrimaryTypeName>.java
type File struct {
	// A Preamble comment belongs not to any type and is usually
	// something like a license or generator header as the first comment In the actual file.
	// The files comment is actually Obj.Comment.
	preamble *Comment
	comment  *Comment
	name     string
	typeDefs []TypeDef
	pkg      *Package
	path     string
	modified bool
}

func (f *File) Modified() bool {
	return f.modified
}

func (f *File) SetModified(modified bool) {
	f.modified = modified
}

func (f *File) Path() string {
	return f.path
}

func (f *File) SetPath(path string) {
	f.path = path
}

func (f *File) Comment() *Comment {
	return f.comment
}

func (f *File) SetComment(comment *Comment) {
	f.comment = comment
}

func (f *File) Pkg() *Package {
	return f.pkg
}

func (f *File) SetPkg(pkg *Package) {
	f.pkg = pkg
}

func NewFile(with func(file *File)) *File {
	f := &File{}
	if with != nil {
		with(f)
	}

	return f
}

func (f *File) Preamble() *Comment {
	return f.preamble
}

func (f *File) SetPreamble(preamble *Comment) {
	f.preamble = preamble
}

func (f *File) Name() string {
	return f.name
}

func (f *File) SetName(name string) {
	f.name = name
}

func (f *File) TypeDefs() []TypeDef {
	return f.typeDefs
}

func (f *File) SetTypeDefs(typeDefs []TypeDef) {
	f.typeDefs = typeDefs
}

func (f *File) AddTypeDefs(typeDefs ...TypeDef) {
	f.typeDefs = append(f.typeDefs, typeDefs...)
}
