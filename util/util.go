package util

import (
	"strings"
	"unicode"
)

func HasPrefixFold(a, b string) bool {
	x := []rune(a)
	y := []rune(b)
	if len(x) < len(y) {
		return false
	}
	for i, c := range y {
		if x[i] == c {
			continue
		}
		if unicode.ToUpper(c) != unicode.ToUpper(x[i]) {
			return false
		}
	}
	return true
}

func SplitCmd(s string) (string, string) {
	firstSpace := strings.Index(s, " ")
	if firstSpace < 0 {
		return s, ""
	}

	return s[:firstSpace], strings.TrimSpace(s[firstSpace:])
}
