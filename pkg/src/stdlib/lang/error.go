package lang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/golang"
	"strings"
	"unicode"
)

// An Error is a sealed (or sum) type of a finite set of enumerable and instantiable types.
// Actually this is just a macro-builder, which creates a macro representing the actual required types.
// per language.
// Go:
//
//	Emits an idiomatic way of representing behavior-based errors instead types. See
//	also https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully and
//	https://dave.cheney.net/2014/12/24/inspecting-errors.
//	Creates private structs each implementing the error interface and methods named after GroupName and each ErrorCase.
//
//	In https://blog.golang.org/error-handling-and-go can also be seen, that the "Is" prefix is omitted (just as the "Get" prefix).
//	Even more verbose (and perhaps unidiomatic), we generate interface types for each error case.
//	We do this only for documentation and reference purpose.
//
// Java:
//
//	model as sealed class or (checked) exception?
type Error struct {
	GroupName string       // GroupName denotes the actual name of the sealed type set of errors.
	Cases     []*ErrorCase // Cases declares possible error cases.
	Comment   string
}

type secretValueErrorKey string

func NewError(groupName string) *Error {
	return &Error{
		GroupName: groupName,
	}
}

// FindError searches in the package of p for an Error group macro with the given name and returns it or nil.
func FindError(p ast.Node, groupName string) *Error {
	pkg := ast.PkgFrom(p)
	if pkg == nil {
		return nil
	}

	key := secretValueErrorKey("error-macro-" + groupName)

	for _, file := range pkg.PkgFiles {
		for _, node := range file.Nodes {
			if macro, ok := node.(*ast.Macro); ok {
				if macro.ID == string(key) {
					if e, ok := macro.Value(key).(*Error); ok {
						return e
					}
				}
			}
		}
	}

	return nil
}

func (n *Error) GetComment() string {
	return n.Comment
}

func (n *Error) Name() string {
	return n.GroupName
}

// SetComment sets the nodes comment.
func (n *Error) SetComment(text string) *Error {
	n.Comment = text
	return n
}

func (n *Error) AddCase(cases ...*ErrorCase) *Error {
	n.Cases = append(n.Cases, cases...)
	for _, errorCase := range cases {
		errorCase.Parent = n
	}

	return n
}

func (n *Error) ID() string {
	return "error-macro-" + n.GroupName
}

// TypeDecl creates a macro to declare the according concrete sealed or sum type(s).
func (n *Error) TypeDecl() *ast.Macro {
	m := ast.NewMacro().SetID(n.ID()).SetMatchers(
		ast.MatchTargetLanguageWithContext(ast.LangGo,
			func(m *ast.Macro) []ast.Node {
				var res []ast.Node

				docSumType := "...returns true, if the error belongs to the sum type of " + n.GroupName + "."
				docUnwrap := "...unpacks the cause or returns nil."
				docErr := "...returns the conventional description of this error."

				sumType := ast.NewInterface(golang.MakePublic(n.GroupName)+"Error").
					SetComment("...represents the sum type behavior of all "+n.GroupName+" errors.").
					AddMethods(
						ast.NewFunc(goErrorMarkerMethod(golang.MakePublic(n.GroupName))).
							SetComment(docSumType).
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))),

						ast.NewFunc("Unwrap").
							SetComment(docUnwrap).
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("error"))),

						ast.NewFunc("Error").
							SetComment(docErr).
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("string"))),
					)

				asSumType := ast.NewFunc("As" + sumType.TypeName).
					SetComment("...finds the first error in err's chain that matches any " + sumType.TypeName + " behavior.\nReturns nil if no such error is found.").
					AddParams(ast.NewParam("err", ast.NewSimpleTypeDecl("error"))).
					AddResults(ast.NewParam("", ast.NewSimpleTypeDecl(ast.Name(sumType.TypeName)))).
					SetBody(ast.NewBlock(
						ast.NewTpl(
							`var match {{.Get "type"}}
								 if {{.Use "errors.As"}}(err, &match) && match.{{.Get "sumTypeMarker"}}() {
									return match
								 }

								 return nil`).
							Put("type", sumType.TypeName).
							Put("sumTypeMarker", goErrorMarkerMethod(golang.MakePublic(n.GroupName))),
					))

				res = append(res, sumType, asSumType)

				for _, errorCase := range n.Cases {
					contract := ast.NewInterface(golang.MakePublic(errorCase.goStructTypeName())).
						SetComment(errorCase.Comment).
						AddEmbedded(ast.NewSimpleTypeDecl(ast.Name(sumType.TypeName)))

					typ := ast.NewStruct(errorCase.goStructTypeName()).
						SetComment(errorCase.Comment + "\n" + errorCase.goStructTypeName() + " is also " + grammarAOrAn(n.GroupName) + "Error.").
						SetVisibility(ast.Private)

					// feed all properties
					for _, property := range errorCase.Properties {
						typ.AddFields(
							ast.NewField(property.goFieldName(), property.decl.Clone()).SetComment(property.comment).SetVisibility(ast.Private),
						)

						// public property getter for struct
						doc := "...returns the value of " + property.name + ".\n" + golang.DeEllipsis(golang.MakePublic(property.name), property.comment)
						typ.AddMethods(
							ast.NewFunc(
								golang.MakePublic(property.goFieldName())).
								SetComment(doc).
								SetRecName("e").
								AddResults(ast.NewParam("", property.decl.Clone())).
								SetBody(
									ast.NewBlock(
										ast.NewTpl("return e." + property.goFieldName()),
									),
								),
						)

						// public property getter for interface
						contract.AddMethods(
							ast.NewFunc(
								golang.MakePublic(property.goFieldName())).
								SetComment(doc).
								AddResults(ast.NewParam("", property.decl.Clone())),
						)

					}

					// insert group marker method
					typ.AddMethods(
						ast.NewFunc(goErrorMarkerMethod(n.GroupName)).
							SetComment(docSumType + "\nThis implementation always returns true.").
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))).
							SetBody(ast.NewBlock(ast.NewReturnStmt(ast.NewIdentLit("true")))),
					)

					// insert case marker method
					doc := "...returns true, if it represents " + grammarAOrAn(errorCase.TypeName) + " case."
					typ.AddMethods(
						ast.NewFunc(goErrorMarkerMethod(errorCase.TypeName)).
							SetComment(doc + "\nThis implementation always returns true.").
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))).
							SetBody(ast.NewBlock(ast.NewReturnStmt(ast.NewIdentLit("true")))),
					)

					contract.AddMethods(
						ast.NewFunc(goErrorMarkerMethod(errorCase.TypeName)).
							SetComment(doc).
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))),
					)

					// always provide an unwrap
					typ.AddFields(ast.NewField("cause", ast.NewSimpleTypeDecl("error")).SetComment("...refers to a causing error or nil.").SetVisibility(ast.Private))
					typ.AddMethods(
						ast.NewFunc("Unwrap").
							SetComment(docUnwrap).
							SetRecName("e").
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("error"))).
							SetBody(
								ast.NewBlock(
									ast.NewTpl("return e.cause"),
								),
							),
					)

					// always provide an error method
					errRetFmt := `return "` + errorCase.Name() + `"`
					if len(errorCase.Properties) > 0 {
						errRetFmt = `return {{.Use "fmt.Sprintf"}}("` + errorCase.Name() + `: `
						args := ""
						for i, property := range errorCase.Properties {
							errRetFmt += property.name + "=%v"
							args += "e." + property.goFieldName()
							if i < len(errorCase.Properties)-1 {
								errRetFmt += ", "
								args += ", "
							}
						}
						errRetFmt += "\", " + args + ")"
					}

					typ.AddMethods(
						ast.NewFunc("Error").
							SetComment(docErr).
							SetRecName("e").
							AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("string"))).
							SetBody(
								ast.NewBlock(
									ast.NewTpl(errRetFmt),
								),
							),
					)

					asType := ast.NewFunc("As" + contract.TypeName).
						SetComment("...finds the first error in err's chain that matches any" + contract.TypeName + " behavior.\nReturns nil if no such error is found.").
						AddParams(ast.NewParam("err", ast.NewSimpleTypeDecl("error"))).
						AddResults(ast.NewParam("", ast.NewSimpleTypeDecl(ast.Name(contract.TypeName)))).
						SetBody(ast.NewBlock(
							ast.NewTpl(
								`var match {{.Get "type"}}
								 if {{.Use "errors.As"}}(err, &match) && match.{{.Get "sumTypeMarker"}}() && match.{{.Get "caseFunc"}}() {
									return match
								 }

								 return nil`).
								Put("type", contract.TypeName).
								Put("caseFunc", goErrorMarkerMethod(errorCase.TypeName)).
								Put("sumTypeMarker", goErrorMarkerMethod(golang.MakePublic(n.GroupName))),
						))

					res = append(res, contract, asType, typ)
				}

				return res
			},
		),
	)

	m.PutValue(secretValueErrorKey(n.ID()), n)

	return m
}

// An ErrorCase declares a unique case of the enumeration.
type ErrorCase struct {
	Parent     *Error
	TypeName   string
	Comment    string
	Properties []errProperty // Properties are usually reflected as Fields and their according getter-method set.
}

func NewErrorCase(name string) *ErrorCase {
	return &ErrorCase{TypeName: name}
}

func (n *ErrorCase) SetComment(text string) *ErrorCase {
	n.Comment = text

	return n
}

func (n *ErrorCase) GetComment() string {
	return n.Comment
}

func (n *ErrorCase) Name() string {
	return n.TypeName
}

func (n *ErrorCase) AddProperty(name string, decl ast.TypeDecl, comment string) *ErrorCase {
	n.Properties = append(n.Properties, errProperty{
		name:    name,
		decl:    decl,
		comment: comment,
	})

	return n
}

// Make creates a new macro, which evaluates to a function or constructor call. Arguments must match exactly
// properties in number and types.
//
//	Go:
//	  - emits a struct literal to a private type, which exposes the according marker interfaces and property getters.
//	Java:
//	  - either throws a (public?) Checked Exception or creates a new instance of a sealed type?
func (n *ErrorCase) Make(args ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguageWithContext(ast.LangGo,
			func(m *ast.Macro) []ast.Node {
				compLit := ast.NewCompLit(ast.NewIdent(n.goStructTypeName()))
				for i, arg := range args {
					compLit.AddElements(ast.NewBinaryExpr(ast.NewIdent(n.Properties[i].name), ast.OpColon, arg))
				}

				return []ast.Node{compLit}
			},
		),
	)
}

type ErrorCheckKind string

const (
	CheckExactBehavior ErrorCheckKind = "exact"
	CheckSumBehavior                  = "sumtype"
	CheckCaseBehavior                 = "singlecase"
)

// Check creates a new macro, which inspects the given variable and tries to match it against this specific case.
//
//	Go:
//	 - creates a new inline type and uses errors.As to Unwrap or match into dstVarName and calls the match block on success.
func (n *ErrorCase) Check(checkKind ErrorCheckKind, checkVarName, dstVarName string, match *ast.Block) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguageWithContext(ast.LangGo,
			func(m *ast.Macro) []ast.Node {
				var iface *ast.Interface
				switch checkKind {
				case CheckExactBehavior:
					iface = n.goFullInterface()
				case CheckSumBehavior:
					iface = n.goSumTypeInterface()
				case CheckCaseBehavior:
					iface = n.goCaseTypeInterface()
				default:
					panic("invalid check kind: " + string(checkKind))
				}

				conditionalType := ""
				for i, fun := range iface.Methods() {
					if len(fun.FunResults) == 1 {
						if sb, ok := fun.FunResults[0].ParamTypeDecl.(*ast.SimpleTypeDecl); ok {
							if sb.Name() == "bool" {
								conditionalType += ` && {{.Get "dst"}}.` + fun.FunName + "()"
							}
						}
					}

					// this is a side effect: in CheckExactBehavior we have at least 2 methods, the first two are always the marker methods,
					// CheckSumBehavior or CheckCaseBehavior have at most 1.
					if i >= 2 {
						break
					}
				}

				return []ast.Node{
					ast.NewTpl(`var {{.Get "dst"}}`).Put("dst", dstVarName),
					iface,
					ast.NewTpl(`
						 if {{.Use "errors.As"}}({{.Get "src"}}, &{{.Get "dst"}})`+conditionalType,
					).Put("dst", dstVarName).Put("src", checkVarName),
					match,
				}
			},
		),
	)
}

func (n *ErrorCase) goSumTypeInterface() *ast.Interface {
	iface := ast.NewInterface("").AddMethods(
		ast.NewFunc(goErrorMarkerMethod(n.Parent.GroupName)).AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))),
	)

	return iface
}

func (n *ErrorCase) goCaseTypeInterface() *ast.Interface {
	iface := ast.NewInterface("").AddMethods(
		ast.NewFunc(goErrorMarkerMethod(n.TypeName)).AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))),
	)

	return iface
}

func (n *ErrorCase) goFullInterface() *ast.Interface {
	iface := ast.NewInterface("").AddMethods(
		ast.NewFunc(goErrorMarkerMethod(n.Parent.GroupName)).AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))),
		ast.NewFunc(goErrorMarkerMethod(n.TypeName)).AddResults(ast.NewParam("", ast.NewSimpleTypeDecl("bool"))),
	)

	for _, property := range n.Properties {
		iface.AddMethods(
			ast.NewFunc(
				golang.MakePublic(property.goFieldName())).
				AddResults(ast.NewParam("", property.decl.Clone())),
		)
	}

	return iface
}

// ContractTypeName returns the (if attached otherwise not-) qualified name of this error type. Depending on the target, a different name (but
// consistent to Error.TypeDecl) referring to an arbitrary type (e.g. interface or class) is returned. If p is not attached to a package,
// just the local name is returned, otherwise a full qualified identifier within the context of p.
func (n *ErrorCase) ContractTypeName(p ast.Node) ast.Name {
	if n.Parent == nil {
		return ast.Name(golang.MakePublic(n.goStructTypeName()))
	}

	var target ast.Target
	mod := &ast.Mod{}
	if ok := ast.ParentAs(p, &mod); ok {
		target = mod.Target
	}

	var identifier string
	switch target.Lang {
	case ast.LangGo:
		identifier = golang.MakePublic(n.goStructTypeName())
	default:
		panic("target lang not yet implemented: " + target.Lang)
	}

	pkg := &ast.Pkg{}
	if ok := ast.ParentAs(p, &pkg); ok {
		return ast.Name(pkg.Path + "." + identifier)
	}

	return ast.Name(identifier)
}

// goStructTypeName is like TicketNotFoundError
func (n *ErrorCase) goStructTypeName() string {
	if n.Parent == nil {
		panic(n.TypeName + " has no error parent")
	}

	const errStr = "Error"
	prefix := n.Parent.GroupName
	if strings.HasSuffix(prefix, errStr) {
		prefix = prefix[:len(prefix)-len(errStr)]
	}

	prefix = golang.MakePrivate(prefix)
	name := prefix + golang.MakePublic(n.TypeName)

	if !strings.HasSuffix(name, errStr) {
		name += errStr
	}

	return name
}

// goErrorMarkerMethod returns a public identifier without any Error suffix.
func goErrorMarkerMethod(s string) string {
	const errStr = "Error"
	if strings.HasSuffix(s, errStr) {
		s = s[:len(s)-len(errStr)]
	}

	return golang.MakePublic(s)
}

type errProperty struct {
	name    string
	decl    ast.TypeDecl
	comment string
}

func (n errProperty) goFieldName() string {
	return golang.MakePrivate(n.name)
}

func grammarAOrAn(s string) string {
	if len(s) == 0 {
		return ""
	}

	switch unicode.ToLower(rune(s[0])) {
	case 'a':
		fallthrough
	case 'e':
		fallthrough
	case 'i':
		return "an " + s
	default:
		return "a " + s
	}

}
