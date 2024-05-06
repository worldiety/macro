package typescript

import (
	"github.com/worldiety/macro/pkg/wdl"
	"unicode"
)

func GetFilename(ident wdl.Identifier) string {
	return string(unicode.ToLower(rune(ident[0]))) + ident.String()[1:] + ".ts"
}
