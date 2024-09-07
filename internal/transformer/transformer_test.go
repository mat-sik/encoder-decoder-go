package transformer

import (
	"bytes"
	"github.com/mat-sik/encoder-decoder/internal/algorithms"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"unicode/utf8"
)

const (
	offset    = 10
	inputRune = '\u2708'
)

var expectedOutputRune = transformFunc(inputRune)
var inputRuneBytes = getInputRuneBytes()

var transformFunc = algorithms.NewOffsetRuneFunc(offset)

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

func Test_applyFuncAndTransfer_properTransfer(t *testing.T) {
	// given
	inputRuneAmount := 4
	reader := new(bytes.Buffer)
	expectedWriter := bytes.NewBuffer(make([]byte, 0, 64))
	for i := 0; i < inputRuneAmount; i++ {
		reader.Write(inputRuneBytes)
		expectedWriter.WriteRune(expectedOutputRune)
	}
	writer := new(bytes.Buffer)

	inputBuffer := bytes.NewBuffer(make([]byte, 0, 2))
	outputBuffer := new(bytes.Buffer)

	expectedInputBufferSize := 0
	expectedOutputBufferSize := 0
	expectedReaderSize := 0

	// when
	err := applyFuncAndTransfer(reader, writer, inputBuffer, outputBuffer, transformFunc)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedInputBufferSize, inputBuffer.Len())
	assert.Equal(t, expectedOutputBufferSize, outputBuffer.Len())
	assert.Equal(t, expectedReaderSize, reader.Len())
	assert.Equal(t, expectedWriter, writer)
}

func Test_applyFuncAndTransfer_lastRuneReadFails(t *testing.T) {
	// given
	reader := new(bytes.Buffer)
	expectedWriter := bytes.NewBuffer(make([]byte, 0, 64))

	reader.Write(inputRuneBytes)
	reader.Write(inputRuneBytes[:1])
	expectedWriter.WriteRune(expectedOutputRune)

	writer := new(bytes.Buffer)

	inputBuffer := bytes.NewBuffer(make([]byte, 0, 2))
	outputBuffer := new(bytes.Buffer)

	expectedInputBufferSize := 1
	expectedOutputBufferSize := 0
	expectedReaderSize := 0

	// when
	err := applyFuncAndTransfer(reader, writer, inputBuffer, outputBuffer, transformFunc)

	// then
	assert.Equal(t, ErrUnableToTransformRune, err)
	assert.Equal(t, expectedInputBufferSize, inputBuffer.Len())
	assert.Equal(t, expectedOutputBufferSize, outputBuffer.Len())
	assert.Equal(t, expectedReaderSize, reader.Len())
	assert.Equal(t, expectedWriter, writer)
}

func Test_applyFuncAndTransfer_consecutiveRuneReadFails(t *testing.T) {
	// given
	reader := new(bytes.Buffer)
	expectedWriter := bytes.NewBuffer(make([]byte, 0, 64))

	reader.Write(inputRuneBytes)
	reader.Write(inputRuneBytes[:1])
	reader.Write(inputRuneBytes)
	expectedWriter.WriteRune(expectedOutputRune)

	writer := new(bytes.Buffer)

	inputBuffer := bytes.NewBuffer(make([]byte, 0, 2))
	outputBuffer := new(bytes.Buffer)

	expectedInputBufferSize := 4
	expectedOutputBufferSize := 0
	expectedReaderSize := 0

	// when
	err := applyFuncAndTransfer(reader, writer, inputBuffer, outputBuffer, transformFunc)

	// then
	assert.Equal(t, ErrUnableToTransformRune, err)
	assert.Equal(t, expectedInputBufferSize, inputBuffer.Len())
	assert.Equal(t, expectedOutputBufferSize, outputBuffer.Len())
	assert.Equal(t, expectedReaderSize, reader.Len())
	assert.Equal(t, expectedWriter, writer)
}

func Test_filesApplyFuncAndTransfer_properTransfer(t *testing.T) {
	// given
	inputFileName := "test_input_file.txt"
	outputFileName := "test_output_file.txt"

	inputFile, err := os.OpenFile(inputFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer closeAndDeleteFile(inputFile)
	runeAmount := 32 * 1024
	initReader(inputFile, runeAmount)

	outputFile, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer closeAndDeleteFile(outputFile)

	inputBuffer := bytes.NewBuffer(make([]byte, 0, ReadBufferSize))
	outputBuffer := bytes.NewBuffer(make([]byte, 0, WriteBufferSize))

	// when
	err = FilesApplyFuncAndTransfer(inputFileName, outputFileName, inputBuffer, outputBuffer, transformFunc)
	// then
	assert.NoError(t, err)
}

func initReader(writer io.Writer, runeAmount int) {
	runesInBufferAmount := 4 * 1024
	bufferSize := runesInBufferAmount * 4
	buffer := make([]byte, 0, bufferSize)
	for i := 0; i < runesInBufferAmount; i++ {
		buffer = utf8.AppendRune(buffer, inputRune)
	}
	for i := 0; i < max(1, runeAmount/runesInBufferAmount); i++ {
		if _, err := writer.Write(buffer); err != nil && err != io.EOF {
			panic(err)
		}
	}
}

func closeAndDeleteFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(file.Name()); err != nil {
		panic(err)
	}
}
