package wdl

import "strings"

func SplitFirstRune(s string) (rune, string) {
	if len(s) == 0 {
		return rune(0), s
	}

	var first rune
	var sb strings.Builder
	for i, r := range s {
		if i == 0 {
			first = r
		} else {
			sb.WriteRune(r)
		}
	}

	return first, sb.String()
}
