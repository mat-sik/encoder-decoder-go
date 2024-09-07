package algorithms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mirrorRuneSlice(t *testing.T) {
	// given
	input := []rune("Hello")
	expected := []rune("Hello")
	// when
	mirrorSlice(input, GetMirrorRuneLatin1)
	mirrorSlice(input, GetMirrorRuneLatin1)
	// then
	assert.Equal(t, expected, input)
}
