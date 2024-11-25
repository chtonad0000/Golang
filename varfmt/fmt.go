//go:build !solution

package varfmt

import (
	"strconv"
	"strings"
)

func interfaceToString(arg interface{}) string {
	switch v := arg.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func Sprintf(format string, args ...interface{}) string {
	indexCount := -1
	i := 0
	var builder strings.Builder

	for i < len(format) {
		if format[i] == '{' {
			k := i + 1
			for format[k] != '}' {
				k++
			}

			indexCount++
			if i+1 == k {
				builder.WriteString(interfaceToString(args[indexCount]))
			} else {
				ind, err := strconv.Atoi(format[i+1 : k])
				if err != nil {
					panic(err)
				}
				builder.WriteString(interfaceToString(args[ind]))
			}
			i = k
		} else {
			builder.WriteByte(format[i])
		}
		i++
	}

	return builder.String()
}
