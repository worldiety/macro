package golang

import (
	"github.com/worldiety/macro/pkg/src/ast"
	"github.com/worldiety/macro/pkg/src/stdlib"
	"strings"
)

// fromStdlib converts stdlib types (indicated by the macro ! sign at the end) and returns a Java name for it.
// Note that primitive types are always returned as their boxed types, because otherwise we would need to carry
// a lot of context information for it. The Java/JVM model is more or less broken for generics and we just wait until
// they fix it up (perhaps with valhalla value types). If you want a reasonable memory usage, you probably
// want a different language anyway.
func fromStdlib(name ast.Name) ast.Name {
	switch name {
	case stdlib.Int:
		return "int"

	case stdlib.Byte:
		return "byte"

	case stdlib.Int16:
		return "int16"

	case stdlib.Int32:
		return "int32"

	case stdlib.Int64:
		return "int64"

	case stdlib.Float32:
		return "float32"

	case stdlib.Float64:
		return "float64"

	case stdlib.Map:
		return "map"

	case stdlib.List:
		return "[]"

	case stdlib.UUID:
		return "github.com/golangee/uuid.UUID"

	case stdlib.String:
		return "string"

	case stdlib.Error:
		return "error"

	case stdlib.Time:
		return "time.Time"

	case stdlib.Duration:
		return "time.Duration"

	case stdlib.URL:
		return "net/url.URL"

	case stdlib.Rune:
		return "rune"

	case stdlib.Void:
		return ""

	case stdlib.Bool:
		return "bool"
	case stdlib.Any:
		return "any"
	default:
		if strings.HasSuffix(string(name), "!") {
			panic("not a stdlib type: " + string(name))
		}
		return name
	}
}
