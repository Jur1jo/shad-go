//go:build !solution

package reverse

import (
	"strings"
)

const (
	oneByteRune = 7 - iota
	tmpRune
	twoByteRune
	thirdByteRune
	fourByteRune
)

const invalidRune = rune('\uFFFD')

func typeOfByte(value byte) int {
	for i := 7; i >= 0; i -= 1 {
		if (value & (1 << i)) == 0 {
			return i
		}
	}
	return -1
}

func Reverse(input string) string {
	var ans strings.Builder
	ans.Grow(len(input))
	currentRune := [...]byte{0, 0, 0, 0}
	indRune := 0
	for i := len(input) - 1; i >= 0; i -= 1 {
		currentRune[indRune] = input[i]
		indRune++
		typeByte := typeOfByte(input[i])
		if typeByte != tmpRune && typeByte >= fourByteRune {
			validRune := false
			switch typeByte {
			case oneByteRune:
				validRune = indRune == 1
			case twoByteRune:
				validRune = indRune == 2
			case thirdByteRune:
				validRune = indRune == 3
			case fourByteRune:
				validRune = indRune == 4
			}
			if validRune {
				for i := indRune - 1; i >= 0; i -= 1 {
					ans.WriteByte(currentRune[i])
				}
				indRune = 0
			} else {
				for i := 0; i < indRune; i++ {
					ans.WriteRune(invalidRune)
				}
				indRune = 0
			}
		} else if typeByte == tmpRune {
			if indRune == 4 {
				ans.WriteRune(invalidRune)
				for i := 3; i > 0; i -= 1 {
					currentRune[i-1] = currentRune[i]
				}
			}
		} else {
			for i := 0; i < indRune; i++ {
				ans.WriteRune(invalidRune)
			}
			indRune = 0
		}
	}
	return ans.String()
}
