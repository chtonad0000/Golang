//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))
	addSpace := false
	for i := 0; i < len(input); {
		run, length := utf8.DecodeRuneInString(input[i:])
		if unicode.IsSpace(run) {
			if !addSpace {
				builder.WriteByte(' ')
				addSpace = true
			}
		} else {
			addSpace = false
			builder.WriteRune(run)
		}
		i += length
	}
	return builder.String()
}
