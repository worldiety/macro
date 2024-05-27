package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
)

func GetFilename(ident wdl.Identifier) string {
	return tsLowerNameStr(ident.String()) + ".ts"
}
