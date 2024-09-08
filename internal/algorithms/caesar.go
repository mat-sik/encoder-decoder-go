package algorithms

import "unicode"

func caesarSlice(input []rune, caesarFunc func(rune) rune) {
	for i, r := range input {
		input[i] = caesarFunc(r)
	}
}

func NewOffsetRuneFunc(offset int32) func(rune) rune {
	if offset > 0 {
		return func(r rune) rune {
			return offsetRuneForward(r, offset)
		}
	} else if offset < 0 {
		return func(r rune) rune {
			return offsetRuneBackward(r, -offset)
		}
	}
	panic("zero offset does not do anything")
}

func offsetRuneForward(r rune, offset int32) rune {
	distance := unicode.MaxRune - r
	if distance < offset {
		return offset - distance - 1
	}
	return r + offset
}

func offsetRuneBackward(r rune, offset int32) rune {
	distance := r
	if distance < offset {
		return unicode.MaxRune - (offset - distance - 1)
	}
	return r - offset
}
