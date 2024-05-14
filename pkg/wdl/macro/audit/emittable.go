package audit

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
	"strconv"
)

const (
	emit = "go.permission.generateTable"
)

type goPermissionTableOptions struct {
	FactoryName    wdl.Identifier `json:"factoryName"`
	PermissionName wdl.Identifier `json:"interfaceName"`
}

type GenerateTable struct {
	prog     *wdl.Program
	preamble string
}

func NewGenerateTable(prog *wdl.Program, preamble string) *GenerateTable {
	return &GenerateTable{prog: prog, preamble: preamble}
}

func (m *GenerateTable) Names() []wdl.MacroName {
	return []wdl.MacroName{
		emit,
	}
}

func (m *GenerateTable) Expand(def wdl.TypeDef, macroInvoc *wdl.MacroInvocation) error {
	var opts goPermissionTableOptions
	if err := macroInvoc.UnmarshalParams(&opts); err != nil {
		return fmt.Errorf("invalid macro params: %w", err)
	}

	if opts.FactoryName == "" {
		opts.FactoryName = "Permissions"
	}

	if opts.PermissionName == "" {
		opts.PermissionName = "Permission"
	}

	def.Pkg().AddFiles(wdl.NewFile(func(file *wdl.File) {
		file.SetMimeType(wdl.MimeTypeGo)

		file.SetName("permissions.gen.go")
		file.SetPath(def.Pkg().Files()[0].Path())
		file.SetModified(true)
		file.SetGenerated(true)
		file.SetPreamble(wdl.NewComment(func(comment *wdl.Comment) {
			comment.AddLines(wdl.NewCommentLine(func(line *wdl.CommentLine) {
				line.SetText(m.preamble)
			}))
		}))

		rec := wdl.NewStruct(func(rec *wdl.Struct) {
			rec.SetName(opts.PermissionName)
			rec.SetPkg(def.Pkg())
			rec.SetComment(wdl.NewSimpleComment(opts.PermissionName.String() + " represents a permission to call a distinct use case. It provides method accessors,\nso that other permission consumers can accept their own interfaces."))

			rec.SetVisibility(wdl.Public)
			rec.AddFields(
				wdl.NewField(func(field *wdl.Field) {
					field.SetName("id")
					field.SetTypeDef(m.prog.MustResolveSimple("std", "string"))
				}),
				wdl.NewField(func(field *wdl.Field) {
					field.SetName("name")
					field.SetTypeDef(m.prog.MustResolveSimple("std", "string"))
				}),
				wdl.NewField(func(field *wdl.Field) {
					field.SetName("desc")
					field.SetTypeDef(m.prog.MustResolveSimple("std", "string"))
				}),
			)

			// add getters for each field, so that other can accept interfaces
			for _, field := range rec.Fields() {
				rec.AddMethods(
					wdl.NewFunc(func(fn *wdl.Func) {
						fn.SetName(field.Name())
						fn.SetVisibility(wdl.Public)
						fn.SetReceiver(wdl.NewParam(func(param *wdl.Param) {
							param.SetName("p")
							param.SetTypeDef(rec.AsResolvedType())
						}))
						fn.SetBody(wdl.NewBlockStmt(func(block *wdl.BlockStmt) {
							block.Add(wdl.RawStmt(fmt.Sprintf("return p.%s", field.Name())))
						}))
						fn.AddResults(wdl.NewParam(func(param *wdl.Param) {
							param.SetTypeDef(m.prog.MustResolveSimple("std", "string"))
						}))
					}),
				)
			}

		})

		fn := wdl.NewFunc(func(fn *wdl.Func) {
			fn.SetName(opts.FactoryName)
			fn.SetComment(wdl.NewSimpleComment(opts.FactoryName.String() + " provides a complete slice of all annotated permissions of all bounded contexts.\nEach use case, which requires some sort of auditing, has its individual permission."))
			fn.SetVisibility(wdl.Public)
			fn.AddResults(wdl.NewParam(func(param *wdl.Param) {
				slice := m.prog.MustResolveSimple("std", "Slice")
				slice.AddParams(rec.AsResolvedType())
				param.SetTypeDef(slice)
			}))
			fn.SetBody(wdl.NewBlockStmt(func(block *wdl.BlockStmt) {
				tmp := fmt.Sprintf("return []%s{\n", rec.Name())
				for _, annotation := range m.prog.Annotations() {
					if pAn, ok := annotation.(*wdl.PermissionAnnotation); ok {
						usecaseFn := pAn.TypeDef().(*wdl.Func)
						ucName := usecaseFn.Name().String()
						if ucAn := wdl.AnnotationForType[*wdl.UseCaseAnnotation](m.prog, usecaseFn); ucAn != nil {
							ucName = ucAn.Name()
						}
						tmp += fmt.Sprintf("{%s, %s, %s},\n", strconv.Quote(pAn.PermissionID()), strconv.Quote(ucName), strconv.Quote(usecaseFn.Comment().String()))
					}
				}
				tmp += "}\n"

				block.Add(wdl.RawStmt(tmp))
			}))
		})

		def.Pkg().AddTypeDefs(fn, rec)
		file.AddTypeDefs(fn, rec)

	}))
	return nil
}
