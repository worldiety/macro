package golang

import (
	"fmt"
	"go/format"
	"strings"
	"unicode"
)

// Format tries to apply the gofmt rules to the given text. If it fails
// the error is returned and the string contains the text with line enumeration.
func Format(source []byte) ([]byte, error) {
	buf, err := format.Source(source)
	if err != nil {
		return []byte(WithLineNumbers(string(source))), err
	}

	return buf, nil
}

// MakePrivate converts ABc to aBc.
// Special cases:
//   - ID becomes id
func MakePrivate(str string) string {
	if len(str) == 0 {
		return str
	}

	switch str {
	case "ID":
		return "id"
	default:
		return string(unicode.ToLower(rune(str[0]))) + str[1:]
	}
}

// MakePublic converts aBc to ABc.
// Special cases:
//   - id becomes ID
func MakePublic(str string) string {
	if len(str) == 0 {
		return str
	}

	switch str {
	case "id":
		return "ID"
	default:
		return string(unicode.ToUpper(rune(str[0]))) + str[1:]
	}
}

// MakeIdentifier creates a public name out of the given string. If it just contains rubbish, at worst the empty
// name _ is returned. - and _ are turned into upper case letters if possible.
func MakeIdentifier(str string) string {
	sb := &strings.Builder{}
	nextUp := true
	first := true
	for _, r := range str {
		if r == '-' || r == '_' || r == ' ' {
			nextUp = true
			continue
		}

		if first && r >= '0' && r <= '9' {
			nextUp = true
			continue
		}

		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			nextUp = true
			continue
		}

		first = false
		if nextUp {
			sb.WriteRune(unicode.ToUpper(r))
			nextUp = false
		} else {
			sb.WriteRune(r)
		}
	}

	if sb.Len() == 0 {
		return "_"
	}

	return sb.String()
}

// WithLineNumbers puts a 1 based line number to the left and returns the text.
func WithLineNumbers(text string) string {
	sb := &strings.Builder{}
	for i, line := range strings.Split(text, "\n") {
		sb.WriteString(fmt.Sprintf("%4d: %s\n", i+1, line))
	}
	return sb.String()
}
