package transformer

import (
	"bytes"
	"fmt"
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

func Test_runeBuffersApplyFuncAndTransfer_allRunesCorrectlyRead(t *testing.T) {
	// given
	inputBuffer := bytes.NewBuffer(make([]byte, 0, 64))
	inputBuffer.WriteRune(inputRune)
	inputBuffer.WriteRune(inputRune)

	outputBuffer := new(bytes.Buffer)

	expectedInputBuffer := bytes.NewBuffer(make([]byte, 0, 64))

	expectedOutputBuffer := new(bytes.Buffer)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	expectedOutputBuffer.WriteRune(expectedOutputRune)

	// when
	err := runeBuffersApplyFuncAndTransfer(inputBuffer, outputBuffer, transformFunc)
	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func Test_runeBuffersApplyFuncAndTransfer_partialTransformationLastRuneNotComplete(t *testing.T) {
	// given
	inputBuffer := new(bytes.Buffer)
	inputBuffer.WriteRune(inputRune)

	inputRuneBytesSlice := inputRuneBytes[:2]
	inputBuffer.Write(inputRuneBytesSlice)

	outputBuffer := new(bytes.Buffer)

	expectedErr := ErrErroneousRune

	expectedInputBuffer := new(bytes.Buffer)
	expectedInputBuffer.Write(inputRuneBytesSlice)

	expectedOutputBuffer := new(bytes.Buffer)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	// when
	err := runeBuffersApplyFuncAndTransfer(inputBuffer, outputBuffer, transformFunc)
	// then
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func Test_runeBuffersApplyFuncAndTransfer_partialTransformationLastRuneNotCompleteAddedMissingRunePart(t *testing.T) {
	// given
	inputBuffer := new(bytes.Buffer)
	inputBuffer.WriteRune(inputRune)

	inputRuneBytesSlice := inputRuneBytes[:2]
	inputBuffer.Write(inputRuneBytesSlice)

	outputBuffer := new(bytes.Buffer)

	expectedErr := ErrErroneousRune

	expectedInputBuffer := bytes.NewBuffer(make([]byte, 0, 64))

	expectedOutputBuffer := new(bytes.Buffer)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	expectedOutputBuffer.WriteRune(expectedOutputRune)
	// when & then
	err := runeBuffersApplyFuncAndTransfer(inputBuffer, outputBuffer, transformFunc)
	assert.Equal(t, expectedErr, err)
	inputBuffer.WriteByte(136)
	err = runeBuffersApplyFuncAndTransfer(inputBuffer, outputBuffer, transformFunc)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func Test_runeBuffersApplyFuncAndTransfer_initialRuneIncomplete(t *testing.T) {
	// given
	inputBuffer := new(bytes.Buffer)
	inputRuneBytesSlice := inputRuneBytes[:2]
	inputBuffer.Write(inputRuneBytesSlice)

	outputBuffer := new(bytes.Buffer)

	expectedErr := ErrErroneousInitialRune

	expectedInputBuffer := new(bytes.Buffer)
	expectedInputBuffer.Write(inputRuneBytesSlice)

	expectedOutputBuffer := new(bytes.Buffer)
	// when
	err := runeBuffersApplyFuncAndTransfer(inputBuffer, outputBuffer, transformFunc)
	// then
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedInputBuffer, inputBuffer)
	assert.Equal(t, expectedOutputBuffer, outputBuffer)
}

func Test_applyFuncAndTransfer_init(t *testing.T) {
	// given
	reader := new(bytes.Buffer)
	for i := 0; i < 4; i++ {
		reader.Write(inputRuneBytes)
	}
	writer := new(bytes.Buffer)

	inputBuffer := bytes.NewBuffer(make([]byte, 0, 2))
	outputBuffer := new(bytes.Buffer)

	// when
	err := applyFuncAndTransfer(reader, writer, inputBuffer, outputBuffer, transformFunc)

	// then
	fmt.Println(err)
}
