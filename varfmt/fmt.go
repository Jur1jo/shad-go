//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Sprintf(format string, args ...interface{}) string {
	var ans strings.Builder
	ans.Grow(len(format))
	lstPos := -1
	countLessFormat := 0
	for i := 0; i < len(format); {
		r, size := utf8.DecodeRuneInString(format[i:])
		switch r {
		case '{':
			if lstPos != -1 {
				panic("Don't correct bracket sequence")
			}
			lstPos = i
		case '}':
			if lstPos == -1 {
				panic("Don't correct bracket sequence")
			}
			if lstPos+1 == i {
				ans.WriteString(fmt.Sprintf("%v", args[countLessFormat]))
			} else {
				arg, _ := strconv.Atoi(format[lstPos+1 : i])
				ans.WriteString(fmt.Sprintf("%v", args[arg]))
			}
			lstPos = -1
			countLessFormat++
		default:
			if lstPos == -1 {
				ans.WriteRune(r)
			}
		}
		i += size
	}
	return ans.String()
}
