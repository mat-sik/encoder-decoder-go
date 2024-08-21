package cipher

import (
	"math"
	"unicode"
)

func mirrorSlice[T comparable](input []T, mirrorFunc func(T) T) {
	for i, r := range input {
		input[i] = mirrorFunc(r)
	}
}

func getMirrorRuneLatin1(r rune) rune {
	if r > unicode.MaxLatin1 {
		panic("incorrect rune provided, at most can be 255")
	}
	return unicode.MaxLatin1 - r
}

func getMirrorRune(r rune) rune {
	return unicode.MaxRune - r
}

func getMirrorByte(b byte) byte {
	return math.MaxUint8 - b
}
