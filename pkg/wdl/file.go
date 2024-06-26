package wdl

import "path/filepath"

type MimeType string

const (
	MimeTypeGo         MimeType = "text/x-go-source"
	MimeTypeTypeScript MimeType = "text/x-typescript-source"
	Raw                MimeType = "binary"
)

// File represents a physical source code file respective compilation unit.
//   - Go: <lowercase AnnotationName>.go
//   - Java: <CamelCasePrimaryTypeName>.java
type File struct {
	// A Preamble comment belongs not to any type and is usually
	// something like a license or generator header as the first comment In the actual file.
	// The files comment is actually Obj.Comment.
	preamble  *Comment
	comment   *Comment
	name      string
	typeDefs  []TypeDef
	pkg       *Package
	path      string
	modified  bool
	generated bool
	imports   map[Identifier]PkgImportQualifier
	mimeType  MimeType
	rawBytes  []byte
}

func (f *File) RawBytes() []byte {
	return f.rawBytes
}

func (f *File) SetRawBytes(rawBytes []byte) {
	f.rawBytes = rawBytes
}

func (f *File) MimeType() MimeType {
	return f.mimeType
}

func (f *File) SetMimeType(mimeType MimeType) {
	f.mimeType = mimeType
}

func (f *File) Import(src *File) {
	f.preamble = src.preamble
	src.preamble = nil

	f.comment = src.comment
	src.comment = nil

	for _, def := range src.typeDefs {
		f.typeDefs = append(f.typeDefs, def)
	}
	src.typeDefs = nil

	f.pkg = src.pkg
	f.path = src.path
	f.modified = true
	f.generated = src.generated

	for identifier, qualifier := range src.imports {
		f.imports[identifier] = qualifier
	}
	clear(src.imports)

}

func (f *File) Generated() bool {
	return f.generated
}

func (f *File) SetGenerated(generated bool) {
	f.generated = generated
}

func (f *File) AbsolutePath() string {
	return filepath.Join(f.path, f.name)
}

func (f *File) Imports() map[Identifier]PkgImportQualifier {
	return f.imports
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
	f := &File{imports: map[Identifier]PkgImportQualifier{}}
	if with != nil {
		with(f)
	}

	return f
}

func (f *File) AddImport(identifier Identifier, qualifier PkgImportQualifier) {
	f.imports[identifier] = qualifier
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
