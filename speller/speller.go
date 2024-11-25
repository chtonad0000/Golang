//go:build !solution

package speller

import (
	"strings"
)

func threeDigitNumberParse(digitMap map[int64]string, n int64) string {
	result := ""
	if n/100 != 0 {
		result += digitMap[n/100] + " hundred "
	}
	if n%100 != 0 {
		result += digitMap[n%100]
	}
	result = strings.TrimSpace(result)
	return result
}

func Spell(n int64) string {
	digitMap := make(map[int64]string)
	result := ""
	ones := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	teens := []string{"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens := []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}

	for i := 0; i < 100; i++ {
		if i < 10 {
			digitMap[int64(i)] = ones[i]
		} else if i < 20 {
			digitMap[int64(i)] = teens[i-10]
		} else {
			tensPart := tens[i/10]
			onesPart := ""
			if i%10 != 0 {
				onesPart = ones[i%10]
				digitMap[int64(i)] = strings.TrimSpace(tensPart + "-" + onesPart)
			} else {
				digitMap[int64(i)] = strings.TrimSpace(tensPart)
			}

		}
	}
	iMap := map[int]string{
		0: "",
		1: "thousand ",
		2: "million ",
		3: "billion ",
	}

	if n == 0 {
		return "zero"
	}
	flag := false
	if n < 0 {
		flag = true
		n = -n
	}
	for i := 0; i < 4; i++ {
		if threeDigitNumberParse(digitMap, n%1000) != "" {
			result = threeDigitNumberParse(digitMap, n%1000) + " " + iMap[i] + result
		}
		n = n / 1000
		if n == 0 {
			break
		}
	}

	result = strings.TrimSpace(result)
	if flag {
		result = "minus " + result
	}
	return result
}
