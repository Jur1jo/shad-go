//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var ans strings.Builder
	ans.Grow(len(input))
	lstIsSpace := false
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		if unicode.IsSpace(r) {
			if !lstIsSpace {
				lstIsSpace = true
				ans.WriteRune(' ')
			}
		} else {
			lstIsSpace = false
			ans.WriteRune(r)
		}
		i += size
	}
	return ans.String()
}
