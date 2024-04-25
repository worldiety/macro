package java

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src"
	"github.com/worldiety/macro/pkg/src/ast"
	"reflect"
	"strconv"
	"strings"
)

const mimeTypeJava = "text/x-java-source"
const packageJavaDocFile = "package-info"

func writeComment(w *src.BufferedWriter, name, doc string) {
	myDoc := formatComment(name, doc)
	if doc != "" {
		w.Printf(myDoc)
		w.Printf("\n")
	}
}

// renderFile tries to emit the file as java
func renderFile(file *ast.SrcFileNode) ([]byte, error) {
	w := &src.BufferedWriter{}

	writeComment(w, file.PkgNode().SrcPackage().PackageName(), file.SrcFile().DocPreamble())
	w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	writeComment(w, file.PkgNode().SrcPackage().PackageName(), file.SrcFile().Doc())

	w.Printf("package %s;\n", file.PkgNode().SrcPackage().PackageName())

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	tmp := &src.BufferedWriter{}
	for _, typ := range file.Types() {
		if err := renderType(typ, tmp); err != nil {
			return nil, err
		}
	}

	// ugly: we may have source file level functions, which are impossible in Java,
	// so we create a new class named <filename>Functions use it to render the functions.
	// This will be package-private only. An alternative would be to pull this up into
	if len(file.SrcFile().Functions()) > 0 {
		artificialHolder := src.NewStruct(file.SrcFile().Name() + "Functions").
			SetDoc("...is introduced to hold static utility functions.").
			SetVisibility(src.PackagePrivate)

		// private constructor
		artificialHolder.AddMethods(
			src.NewFunc(artificialHolder.Name()).
				SetVisibility(src.Private).
				SetDoc("...is a private constructor because this class only contains static methods.").
				SetBody(src.NewBlock()),
		)

		artificialHolder.AddMethods(file.SrcFile().Functions()...)
		node := ast.NewTypeNode(file, artificialHolder)

		if err := renderType(node, tmp); err != nil {
			return nil, err
		}
	}

	importer := importerFromTree(file)
	for _, qualifier := range importer.qualifiers() {
		w.Printf("import %s;\n", qualifier)
	}

	w.Printf(tmp.String())

	return Format(w.Bytes())
}

func renderTypePreamble(w *src.BufferedWriter, node interface {
	Name() string
	Doc() string
	Annotations() []*ast.AnnotationNode
}) error {
	writeComment(w, node.Name(), node.Doc())

	for _, annotation := range node.Annotations() {
		if err := renderAnnotation(annotation, w); err != nil {
			return err
		}
		w.Printf("\n")
	}

	return nil
}

func renderType(t *ast.TypeNode, w *src.BufferedWriter) error {

	switch node := t.NamedNode().(type) {
	case *ast.StructNode:
		return renderStruct(node, w)
	case *ast.InterfaceNode:
		return renderInterface(node, w)
	default:
		panic("type not yet implemented: " + reflect.TypeOf(t).String())
	}
}

func renderInterface(node *ast.InterfaceNode, w *src.BufferedWriter) error {
	if err := renderTypePreamble(w, node); err != nil {
		return err
	}

	w.Printf(visibilityAsKeyword(node.SrcInterface().Visibility()))

	w.Printf(" interface %s {\n", node.SrcInterface().Name())

	for _, typeNode := range node.Types() {
		if err := renderType(typeNode, w); err != nil {
			return err
		}
	}

	for _, fun := range node.Methods() {
		if err := renderFunc(fun, w); err != nil {
			return fmt.Errorf("failed to render func %s: %w", fun.SrcFunc().Name(), err)
		}
	}
	w.Printf("}\n")

	return nil
}

func renderStruct(node *ast.StructNode, w *src.BufferedWriter) error {
	if err := renderTypePreamble(w, node); err != nil {
		return err
	}

	w.Printf(visibilityAsKeyword(node.SrcStruct().Visibility()))
	if node.SrcStruct().Final() {
		w.Printf(" final ")
	}

	if node.SrcStruct().Static() {
		w.Printf(" static ")
	}

	w.Printf(" class %s {\n", node.SrcStruct().Name())

	for _, typeNode := range node.Types() {
		if err := renderType(typeNode, w); err != nil {
			return err
		}
	}

	for _, field := range node.Fields() {
		if err := renderField(field, w); err != nil {
			return fmt.Errorf("failed to render field %s: %w", field.SrcField().Name(), err)
		}
	}

	for _, fun := range node.Methods() {
		if err := renderFunc(fun, w); err != nil {
			return fmt.Errorf("failed to render func %s: %w", fun.SrcFunc().Name(), err)
		}
	}
	w.Printf("}\n")

	return nil
}

func renderFunc(node *ast.FuncNode, w *src.BufferedWriter) error {
	comment := &strings.Builder{}
	comment.WriteString(node.SrcFunc().Doc())
	comment.WriteString("\n\n")

	for _, parameterNode := range node.InputParams() {
		if parameterNode.SrcParameter().Doc() == "" {
			continue
		}

		comment.WriteString("@param ")
		name := parameterNode.SrcParameter().Name()
		if name == "" {
			name = fromStdlib(src.Name(parameterNode.SrcParameter().TypeDecl().String())).Identifier()
		}

		comment.WriteString(deEllipsis(name, parameterNode.SrcParameter().Doc()))
		comment.WriteString("\n")
	}

	for i, parameterNode := range node.OutputParams() {
		if i == 0 || parameterNode.SrcParameter().Doc() == "" {
			continue
		}

		comment.WriteString("@throws ")
		name := parameterNode.SrcParameter().Name()
		if name == "" {
			name = fromStdlib(src.Name(parameterNode.SrcParameter().TypeDecl().String())).Identifier()
		}

		comment.WriteString(deEllipsis(name, parameterNode.SrcParameter().Doc()))
		comment.WriteString("\n")
	}

	writeComment(w, node.SrcFunc().Name(), comment.String())

	if _, ok := node.Parent().(*ast.InterfaceNode); ok {
		// we ignore the visibility entirely, because in Java interfaces methods are always public
	} else {
		w.Printf(visibilityAsKeyword(node.SrcFunc().Visibility()))
		w.Printf(" ")
		if node.SrcFunc().Static() {
			w.Printf("static ")
		}
	}

	for _, annotation := range node.Annotations() {
		if err := renderAnnotation(annotation, w); err != nil {
			return err
		}
		w.Printf("\n")
	}

	if len(node.OutputParams()) == 0 {
		parentType := ast.ParentTypeNode(node)
		if node.SrcFunc().Name() == parentType.SrcNamedType().Name() {
			// special case, if we are a constructor, we omit also the void
		} else {
			w.Printf("void ")
		}
	} else {
		if err := renderTypeDecl(node.OutputParams()[0].TypeDecl(), w); err != nil {
			return err
		}
		w.Printf(" ")
	}
	w.Printf(node.SrcFunc().Name())
	w.Printf("(")
	for i, parameterNode := range node.InputParams() {
		for _, annotationNode := range parameterNode.Annotations() {
			if err := renderAnnotation(annotationNode, w); err != nil {
				return err
			}

			w.Printf(" ")
		}

		if err := renderTypeDecl(parameterNode.TypeDecl(), w); err != nil {
			return err
		}

		if i == len(node.InputParams())-1 && node.SrcFunc().Variadic() {
			w.Printf("...")
		} else {
			w.Printf(" ")
		}

		w.Printf(parameterNode.SrcParameter().Name())

		if i < len(node.OutputParams())-1 {
			w.Printf(", ")
		}
	}
	w.Printf(")")

	// by convention this must be throwables in Java
	if len(node.OutputParams()) > 1 {
		w.Printf("throws ")
		for i, parameterNode := range node.OutputParams() {
			if i == 0 {
				continue
			}

			if err := renderTypeDecl(parameterNode.TypeDecl(), w); err != nil {
				return err
			}

			if i < len(node.OutputParams())-1 {
				w.Printf(", ")
			}
		}
	}

	if node.SrcFunc().Body() == nil {
		w.Printf(";")
	} else {
		w.Printf("{\n")
		w.Printf("}\n")
	}

	return nil
}

func renderField(node *ast.FieldNode, w *src.BufferedWriter) error {
	writeComment(w, node.SrcField().Name(), node.SrcField().Doc())
	for _, annotation := range node.Annotations() {
		if err := renderAnnotation(annotation, w); err != nil {
			return err
		}
	}
	w.Printf(visibilityAsKeyword(node.SrcField().Visibility()))
	w.Printf(" ")
	if err := renderTypeDecl(node.TypeDecl(), w); err != nil {
		return err
	}
	w.Printf(" ")
	w.Printf(node.SrcField().Name())
	w.Printf(";\n")

	return nil
}

func renderAnnotation(node *ast.AnnotationNode, w *src.BufferedWriter) error {
	importer := importerFromTree(node)

	w.Printf("@")
	w.Printf(string(importer.shortify(node.SrcAnnotation().Name())))
	attrs := node.SrcAnnotation().Attributes()
	if len(attrs) > 0 {
		w.Printf("(")
		// the default case
		if len(attrs) == 1 && attrs[0] == "" {
			w.Printf(node.SrcAnnotation().Value(""))
		} else {
			// the named attribute cases
			for i, attr := range attrs {
				w.Printf(attr)
				w.Printf(" = ")
				w.Printf(node.SrcAnnotation().Value(attr))
				if i < len(attrs)-1 {
					w.Printf(", ")
				}
			}
		}

		w.Printf(")")
	}

	return nil
}

func renderTypeDecl(node ast.TypeDeclNode, w *src.BufferedWriter) error {
	importer := importerFromTree(node)

	switch t := node.(type) {
	case *ast.SimpleTypeDeclNode:
		w.Printf(string(importer.shortify(fromStdlib(t.SrcSimpleTypeDecl().Name()))))
	case *ast.TypeDeclPtrNode:
		atomicReference := importer.shortify("java.util.concurrent.atomic.AtomicReference")
		w.Printf(string(atomicReference) + "<")
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf(">")
	case *ast.SliceTypeDeclNode:
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf("[]")
	case *ast.GenericTypeDeclNode:
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf("<")
		for i, decl := range t.Params() {
			if err := renderTypeDecl(decl, w); err != nil {
				return err
			}
			if i < len(t.Params())-1 {
				w.Printf(",")
			}
		}
		w.Printf(">")
	case *ast.ChanTypeDeclNode:
		blockingQueue := importer.shortify("java.util.concurrent.BlockingQueue")
		w.Printf(string(blockingQueue) + "<")
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf(">")

	case *ast.ArrayTypeDeclNode:
		// in Java this is the same as a slice, we cannot have yet custom size value arrays. Perhaps
		// valhalla may fix that
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf("[]")
	case *ast.FuncTypeDeclNode:
		// Java does not have it. We would need to create a functional interface for it, which is out of scope here.
		writeComment(w, "", "inline function declarations are not supported by java:\n\n"+t.SrcFuncTypeDecl().String())
		w.Printf("Object")
		return nil
	default:
		panic("not yet implemented: " + reflect.TypeOf(t).String())
	}

	return nil
}

func visibilityAsKeyword(v src.Visibility) string {
	switch v {
	case src.Public:
		return "public"
	case src.PackagePrivate:
		return ""
	case src.Private:
		return "private"
	case src.Protected:
		return "protected"
	default:
		panic("visibility not implemented: " + strconv.Itoa(int(v)))
	}

}
