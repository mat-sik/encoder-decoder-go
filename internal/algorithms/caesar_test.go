package algorithms

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode"
)

func Test_caesarSlice(t *testing.T) {
	// given
	input := []rune("aaaaa")
	expected := []rune("aaaaa")
	var offset int32 = 1
	// when
	caesarSlice(input, offset, OffsetRuneForward)
	fmt.Println(string(input))
	caesarSlice(input, offset, OffsetRuneBackward)
	fmt.Println(string(input))
	// then
	assert.Equal(t, expected, input)
}

func Test_offsetRuneForward(t *testing.T) {
	// given
	var input rune = 97
	fmt.Println(input)
	fmt.Println(unicode.MaxRune)
	// when
	var output rune = OffsetRuneForward(input, 1)
	fmt.Println(output)
	fmt.Println(string(output))
	// then
}

func Test_offsetRuneBackward(t *testing.T) {
	// given
	var input rune = 5
	fmt.Println(input)
	fmt.Println(unicode.MaxRune)
	// when
	var output rune = OffsetRuneBackward(input, 4)
	fmt.Println(output)
	fmt.Println(string(output))
	// then
}
