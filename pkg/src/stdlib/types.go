package stdlib

const (
	Any = "any!"

	// Bool renders to the bool in Go or to the java boolean.
	Bool = "bool!"

	// Int renders platform specific, for Go this is either signed 32 or 64bit but for Java always signed 32bit.
	Int = "int!"

	// Byte renders to the unsigned byte or the java signed byte.
	Byte = "byte!"

	// Int16 renders to short or int16.
	Int16 = "int16!"

	// Int32 renders to int for Java or else int32
	Int32 = "int32!"

	// Int64 renders to long or int64.
	Int64 = "int64!"

	// Float32 renders as a float or float32.
	Float32 = "float32!"

	// Float64 renders as a double or float64.
	Float64 = "float64!"

	// Map refers to e.g. the Go build-in type map or the java.util.Map<K,V> type.
	Map = "map!"

	// List refers to e.g. the Go build-in slice or the java.util.List<T> type.
	List = "list!"

	// UUID refers to either github.com/golangee/uuid.UUID or java.util.UUID.
	UUID = "uuid!"

	// String refers to the Go build-in string or the java.lang.String type.
	String = "string!"

	// Error refers to the Go build-in error type or java.lang.Exception.
	Error = "error!"

	// Time refers to the Go time.Time type or java.time.ZonedDateTime.
	Time = "time!"

	// Duration refers to the Go time.Duration or java.time.Duration type.
	Duration = "duration!"

	// URL refers to the Go *net/url.URL or java.net.URL type.
	URL = "url!"

	// Rune represents a 32bit unicode codepoint.
	Rune = "rune!"

	// Void represents nothing and should only be used for languages which requires to say that nothing is really
	// nothing.
	Void = "void!"
)

// Types returns all defined standard library transpiler types.
var Types = []string{
	Bool,
	Int,
	Byte,
	Int16,
	Int32,
	Int64,
	Float32,
	Float64,
	Map,
	List,
	UUID,
	String,
	Error,
	Time,
	Duration,
	URL,
	Rune,
	Void,
}

// IsNumber returns true for any int or float.
func IsNumber(n string) bool {
	switch n {
	case Int:
		fallthrough
	case Byte:
		fallthrough
	case Int16:
		fallthrough
	case Int32:
		fallthrough
	case Int64:
		fallthrough
	case Float32:
		fallthrough
	case Float64:
		return true
	}

	return false
}
