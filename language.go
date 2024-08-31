package program

import (
	"unicode"
)

func sentence(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	if runes[len(runes)-1] != '.' {
		runes = append(runes, '.')
	}

	return string(runes)
}
