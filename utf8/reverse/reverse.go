//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))
	i := len(input)
	for i > 0 {
		run, length := utf8.DecodeLastRuneInString(input[:i])
		builder.WriteRune(run)
		i -= length
	}
	return builder.String()
}
