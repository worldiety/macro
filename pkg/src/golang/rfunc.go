package golang

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/golang/validate"
	"github.com/worldiety/macro/pkg/src/render"
	"strings"
)

// renderFunc inspects and emits a function signature with optional receiver and optional body, depending
// on the actual parent, which is either an ast.File, ast.Struct or ast.Interface. The keyword func is not
// rendered here, because we renderFunc also for type declarations.
func (r *Renderer) renderFunc(node *ast.Func, w *render.BufferedWriter) error {
	funComment := r.renderFuncComment(node)
	if funComment != "" {
		r.writeComment(w, false, node.Identifier(), funComment)
	}

	var structNode *ast.Struct
	switch t := node.Parent().(type) {
	case *ast.Struct:
		structNode = t
		recName := "_"
		if node.RecName() != "" {
			recName = node.RecName()
		}

		if node.PtrReceiver() {
			recName += "*"
		}

		w.Printf("func (%s %s) ", recName, t.Identifier())
	case *ast.Interface:

	default:
		w.Printf("func ")
	}

	if err := validate.ExportedIdentifier(node.Visibility(), node.Identifier()); err != nil {
		return err
	}

	w.Printf(node.Identifier())
	w.Printf("(")
	for i, parameterNode := range node.FunParams {

		w.Printf(parameterNode.Identifier())

		if i == len(node.FunParams)-1 && node.Variadic() {
			w.Printf("...")
		} else {
			w.Printf(" ")
		}

		if err := r.renderTypeDecl(parameterNode.TypeDecl(), w); err != nil {
			return fmt.Errorf("unable to render input parameter TypeDecl: %w", err)
		}

		if i < len(node.FunParams)-1 {
			w.Printf(", ")
		}
	}
	w.Printf(")")

	if len(node.FunResults) > 0 {
		w.Printf("(")
	}
	for i, parameterNode := range node.FunResults {
		w.Printf(parameterNode.Identifier())
		w.Printf(" ")

		if err := r.renderTypeDecl(parameterNode.TypeDecl(), w); err != nil {
			return fmt.Errorf("unable to render ouput parameter TypeDecl: %w", err)
		}

		if i < len(node.FunResults)-1 {
			w.Printf(", ")
		}
	}
	if len(node.FunResults) > 0 {
		w.Printf(")")
	}

	if node.Body() == nil {
		if structNode != nil {
			return fmt.Errorf("a struct method must have a body")
		}

		w.Printf("\n")
	} else {
		if err := r.renderNode(node.Body(), w); err != nil {
			return fmt.Errorf("unable to render function body: %w", err)
		}
	}

	return nil
}

func (r *Renderer) renderFuncComment(node *ast.Func) string {
	comment := &strings.Builder{}
	if node.ObjComment != nil {
		comment.WriteString(node.ObjComment.Text)
	}

	hasParamComments := false
	for _, param := range node.FunParams {
		if param.ObjComment != nil {
			hasParamComments = true
			break
		}
	}

	for _, param := range node.Results() {
		if param.ObjComment != nil {
			hasParamComments = true
			break
		}
	}

	if hasParamComments || len(node.ErrorHintRefs) > 0 {
		comment.WriteString("\n\n")
	}

	for _, parameterNode := range node.FunParams {
		if parameterNode.ObjComment == nil {
			continue
		}

		comment.WriteString("The parameter ")
		name := parameterNode.Identifier()
		if name == "" {
			name = fromStdlib(ast.Name(parameterNode.TypeDecl().String())).Identifier()
		}

		comment.WriteString(DeEllipsis(name, parameterNode.Obj.ObjComment.Text))
		comment.WriteString("\n")
	}

	for i, parameterNode := range node.FunResults {
		if i == 0 || parameterNode.ObjComment == nil {
			continue
		}

		comment.WriteString("The result ")
		name := parameterNode.Identifier()
		if name == "" {
			name = fromStdlib(ast.Name(parameterNode.TypeDecl().String())).Identifier()
		}

		comment.WriteString(DeEllipsis(name, parameterNode.ObjComment.Text))
		comment.WriteString("\n")
	}

	for _, ref := range node.ErrorHintRefs {
		details := strings.TrimSpace(DeEllipsis("", ref.GetComment()))
		if details != "" {
			comment.WriteString("Returns error '" + ref.Name() + "' when " + details)
		} else {
			comment.WriteString("Returns error '" + ref.Name() + "'.")
		}
		comment.WriteString("\n")
	}

	return comment.String()
}
