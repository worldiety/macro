package validate

import (
	"fmt"
	"github.com/worldiety/macro/pkg/src/ast"
	"unicode"
)

// Keywords contains all go keywords which may not be used as identifiers.
// See https://golang.org/ref/spec#Keywords.
var Keywords = []string{"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
}

// Identifier asserts the given string is an identifier.
// See https://golang.org/ref/spec#Identifiers.
func Identifier(identifier string) error {
	if identifier == "" {
		return fmt.Errorf("an empty string is not a valid identifier")
	}

	for i, r := range identifier {
		if i == 0 && !(unicode.IsLetter(r) || r == '_') {
			return fmt.Errorf("the first char '%s' of identifier '%s' must be a unicode letter", string(r), identifier)
		} else {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
				return fmt.Errorf("the %d rune ('%s') of identifier '%s' must be a unicode letter or digit", i, string(r), identifier)
			}
		}
	}

	for _, keyword := range Keywords {
		if keyword == identifier {
			return fmt.Errorf("identifier '%s' is a keyword", identifier)
		}
	}

	return nil
}

// ExportedIdentifier asserts the identifier matches the visibility exporting rules.
// See https://golang.org/ref/spec#Exported_identifiers.
func ExportedIdentifier(visibility ast.Visibility, identifier string) error {
	if err := Identifier(identifier); err != nil {
		return err
	}

	if identifier == "_" {
		return nil
	}

	var firstRune rune
	for _, r := range identifier {
		firstRune = r
		break
	}

	switch visibility {
	case ast.Public:
		if !unicode.IsUpper(firstRune) {
			return fmt.Errorf("expected '%s' to be an exported identifier", identifier)
		}
	case ast.PackagePrivate:
		fallthrough
	case ast.Protected:
		fallthrough
	case ast.Private:
		if !unicode.IsLower(firstRune) {
			return fmt.Errorf("expected '%s' to be an unexported (package private) identifier", identifier)
		}
	}

	return nil
}
