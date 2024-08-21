package cipher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mirrorRuneSlice(t *testing.T) {
	// given
	input := []rune("Hello")
	expected := []rune("Hello")
	// when
	mirrorSlice(input, getMirrorRuneLatin1)
	mirrorSlice(input, getMirrorRuneLatin1)
	// then
	assert.Equal(t, expected, input)
}
