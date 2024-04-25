package ast

// Arch defines the architecture to generate code for.
type Arch string

const (
	ArchAMD64 Arch = "amd64"
	ArchARM64 Arch = "arm64"
	ArchWASM  Arch = "wasm"
)

// OS defines the operating system to generate code for.
type OS string

const (
	OSLinux  OS = "linux"
	OSWin    OS = "windows"
	OSDarwin OS = "darwin"
	OSIOS    OS = "iOS"
)

// Lang defines the target language to generate code for.
type Lang string

const (
	LangJava   Lang = "java"
	LangGo     Lang = "go"
	LangRust   Lang = "rust"
	LangSwift  Lang = "swift"
	LangKotlin Lang = "kotlin"
	LangC      Lang = "c"
	LangCPP    Lang = "c++"
)

// LangVersion specifies an arbitrary version string for a specific language. There is no guarantee of a semantic
// version logic here.
type LangVersion string

const (
	LangVersionJava8 LangVersion = "1.8"
	LangVersionGo16  LangVersion = "1.16"
	LangVersionGo17  LangVersion = "1.17"
	LangVersionSwift LangVersion = "5.1"
)

// Framework specifies an arbitrary framework string. This is mostly unique In conjunction with a language.
type Framework string

const (
	FrameworkSDK Framework = ""
)

type Target struct {
	Out            string      // the output directory
	Arch           Arch        // probably empty and of limited use.
	Os             OS          // probably empty and of limited use.
	Lang           Lang        // target generator language.
	MinLangVersion LangVersion // the (inclusive) supported minimum version of the generated code.
	MaxLangVersion LangVersion // the (inclusive) supported maximum version of the generated code.
	Framework      Framework   // the framework to use. Empty means only use the default standard library things.
	Require        struct { // require directive
		GoMod []string // go mod specific directive strings (e.g. github.com/golangee/sql v0.0.0-20210531101020-33021aed64c2)
	}
}

func (t Target) Equals(o Target) bool {
	return t.Lang == o.Lang && t.Os == o.Os && t.Arch == o.Arch && t.MinLangVersion == o.MinLangVersion && t.MaxLangVersion == o.MaxLangVersion && t.Framework == o.Framework
}

// A Mod is the root of a project and describes a module with packages.
//  * Java: denotes a gradle module (build.gradle).
//  * Go: describes a Go module (go.mod).
type Mod struct {
	Name   string // Name refers to a unique module name. In go this is the module name.
	Target Target
	Pkgs   []*Pkg
	Obj
}

// NewMod allocates a new Module.
func NewMod(name string) *Mod {
	return &Mod{Name: name}
}

// SetLang updates the Target.Lang
func (n *Mod) SetLang(lang Lang) *Mod {
	n.Target.Lang = lang
	return n
}

// SetOutputDirectory sets the targets output directory.
func (n *Mod) SetOutputDirectory(dir string) *Mod {
	n.Target.Out = dir
	return n
}

// SetLangVersion updates the Target.MinLangVersion.
func (n *Mod) SetLangVersion(version LangVersion) *Mod {
	n.Target.MinLangVersion = version
	return n
}

// Require expects a language to decide how to handle the dependency.
func (n *Mod) Require(dep string) *Mod {
	if n.Target.Lang != LangGo {
		panic("invalid state: Require currently only supports Go")
	}

	n.Target.Require.GoMod = append(n.Target.Require.GoMod, dep)
	return n
}

// AddPackages appends the given packages and updates the Parent accordingly.
func (n *Mod) AddPackages(packages ...*Pkg) *Mod {
	n.Pkgs = append(n.Pkgs, packages...)
	for _, pkg := range packages {
		assertNotAttached(pkg)
		pkg.Obj.ObjParent = n
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Mod) Children() []Node {
	tmp := make([]Node, 0, len(n.Pkgs)+1)
	tmp = append(tmp, n.Obj.ObjComment)

	for _, pkg := range n.Pkgs {
		tmp = append(tmp, pkg)
	}

	return tmp
}

// Doc sets the nodes comment.
func (n *Mod) Doc(text string) *Mod {
	n.Obj.ObjComment = NewComment(text)
	n.Obj.ObjComment.ObjParent = n
	return n
}
