package transformer

import (
	"bytes"
	"github.com/mat-sik/encoder-decoder/internal/algorithms"
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode/utf8"
)

const (
	offset    = 10
	inputRune = '\u2708'
)

var expectedOutputRune = transformFunc(inputRune)
var inputRuneBytes = getInputRuneBytes()

var transformFunc = func(r rune) rune {
	return algorithms.OffsetRuneForward(r, offset)
}

func getInputRuneBytes() []byte {
	inputRuneBytes := make([]byte, utf8.RuneLen(inputRune))
	utf8.EncodeRune(inputRuneBytes, inputRune)
	return inputRuneBytes
}

func Test_transformRuneBuffers_allRunesCorrectlyRead(t *testing.T) {
	// given
	inputBuffer := new(bytes.Buffer)
	inputBuffer.WriteRune(inputRune)
	inputBuffer.WriteRune(inputRune)

	outputBuffer := new(bytes.Buffer)

	expectedInputBuffer := getEmptyInitialisedBuffer()

	expectedOutputBuffer := new(bytes.Buffer)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	expectedOutputBuffer.WriteRune(expectedOutputRune)

	// when
	err := transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)
	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func Test_transformRuneBuffers_partialTransformationLastRuneNotComplete(t *testing.T) {
	// given
	inputBuffer := new(bytes.Buffer)
	inputBuffer.WriteRune(inputRune)

	inputRuneBytesSlice := inputRuneBytes[:2]
	inputBuffer.Write(inputRuneBytesSlice)

	outputBuffer := new(bytes.Buffer)

	expectedErr := &ErrErroneousRune{}

	expectedInputBuffer := new(bytes.Buffer)
	expectedInputBuffer.Write(inputRuneBytesSlice)

	expectedOutputBuffer := new(bytes.Buffer)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	// when
	err := transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)
	// then
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func Test_transformRuneBuffers_partialTransformationLastRuneNotCompleteAddedMissingRunePart(t *testing.T) {
	// given
	inputBuffer := new(bytes.Buffer)
	inputBuffer.WriteRune(inputRune)

	inputRuneBytesSlice := inputRuneBytes[:2]
	inputBuffer.Write(inputRuneBytesSlice)

	outputBuffer := new(bytes.Buffer)

	expectedErr := &ErrErroneousRune{}

	expectedInputBuffer := getEmptyInitialisedBuffer()

	expectedOutputBuffer := new(bytes.Buffer)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	// when & then
	err := transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)
	assert.Equal(t, expectedErr, err)
	inputBuffer.WriteByte(136)
	err = transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func Test_transformRuneBuffers_initialRuneIncomplete(t *testing.T) {
	// given
	inputBuffer := new(bytes.Buffer)
	inputRuneBytesSlice := inputRuneBytes[:2]
	inputBuffer.Write(inputRuneBytesSlice)

	outputBuffer := new(bytes.Buffer)

	expectedErr := &ErrErroneousInitialRune{}

	expectedInputBuffer := new(bytes.Buffer)
	expectedInputBuffer.Write(inputRuneBytesSlice)

	expectedOutputBuffer := new(bytes.Buffer)
	// when
	err := transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)
	// then
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func getEmptyInitialisedBuffer() *bytes.Buffer {
	buffer := new(bytes.Buffer)
	buffer.Grow(64)
	return buffer
}
