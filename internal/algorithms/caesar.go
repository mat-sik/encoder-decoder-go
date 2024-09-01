package algorithms

import "unicode"

func caesarSlice(input []rune, offset int32, caesarFunc func(rune, int32) rune) {
	for i, r := range input {
		input[i] = caesarFunc(r, offset)
	}
}

func OffsetRuneForward(r rune, offset int32) rune {
	distance := unicode.MaxRune - r
	if distance < offset {
		return offset - distance - 1
	}
	return r + offset
}

func OffsetRuneBackward(r rune, offset int32) rune {
	distance := r
	if distance < offset {
		return unicode.MaxRune - (offset - distance - 1)
	}
	return r - offset
}
