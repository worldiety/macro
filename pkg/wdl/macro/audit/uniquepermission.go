package audit

import (
	"fmt"
	"github.com/worldiety/macro/pkg/wdl"
)

type UniquePermissionCheckAnnotation struct {
	prog *wdl.Program
}

func NewUniquePermissionCheckAnnotation(prog *wdl.Program) *UniquePermissionCheckAnnotation {
	return &UniquePermissionCheckAnnotation{prog: prog}
}

func (m *UniquePermissionCheckAnnotation) Expand() error {
	tmp := map[string]*wdl.PermissionAnnotation{}
	for _, annotation := range m.prog.Annotations() {
		if pAn, ok := annotation.(*wdl.PermissionAnnotation); ok {
			alreadyDefined, ok := tmp[pAn.PermissionID()]
			if ok {
				return fmt.Errorf("duplicate permission annotation found: %s: %s and %s", pAn.PermissionID(), alreadyDefined.TypeDef().Name(), pAn.TypeDef().Name())
			}

			tmp[pAn.PermissionID()] = pAn
		}
	}

	return nil
}
