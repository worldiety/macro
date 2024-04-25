package java

import (
	"github.com/worldiety/macro/pkg/src"
	"github.com/worldiety/macro/pkg/src/stdlib"
	"strings"
)

// fromStdlib converts stdlib types (indicated by the macro ! sign at the end) and returns a Java name for it.
// Note that primitive types are always returned as their boxed types, because otherwise we would need to carry
// a lot of context information for it. The Java/JVM model is more or less broken for generics and we just wait until
// they fix it up (perhaps with valhalla value types). If you want a reasonable memory usage, you probably
// want a different language anyway.
func fromStdlib(name src.Name) src.Name {
	switch name {
	case stdlib.Int:
		return "Integer"

	case stdlib.Byte:
		return "Byte"

	case stdlib.Int16:
		return "Short"

	case stdlib.Int32:
		return "Integer"

	case stdlib.Int64:
		return "Long"

	case stdlib.Float32:
		return "Float"

	case stdlib.Float64:
		return "Double"

	case stdlib.Map:
		return "java.util.Map"

	case stdlib.List:
		return "java.util.List"

	case stdlib.UUID:
		return "java.util.UUID"

	case stdlib.String:
		return "String"

	case stdlib.Error:
		return "Exception"

	case stdlib.Time:
		return "java.time.ZonedDateTime"

	case stdlib.Duration:
		return "java.time.Duration"

	case stdlib.URL:
		return "java.net.URL"

	case stdlib.Rune:
		// in Java this is just an int, using chars is broken and legacy code, see
		// https://docs.oracle.com/javase/tutorial/i18n/text/characterClass.html
		return "int"

	case stdlib.Void:
		return "void"

	default:
		if strings.HasSuffix(string(name), "!") {
			panic("not a stdlib type: " + string(name))
		}
		return name
	}
}
