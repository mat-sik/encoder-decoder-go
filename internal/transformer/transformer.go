package transformer

import (
	"bytes"
	"errors"
	"io"
	"os"
)

const ReadBufferSize = 4 * 1024
const WriteBufferSize = 4 * ReadBufferSize

func transformAndTransferRuneFiles(inputFilePath string, outputFilePath string, transformFunc func(r rune) rune) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer closeFile(inputFile)

	outputFile, err := os.Open(outputFilePath)
	if err != nil {
		return err
	}
	defer closeFile(outputFile)

	inputBuffer := new(bytes.Buffer)
	inputBuffer.Grow(ReadBufferSize)

	outputBuffer := new(bytes.Buffer)
	outputBuffer.Grow(WriteBufferSize)

	if err = transformAndTransferRunes(inputFile, outputFile, inputBuffer, outputBuffer, transformFunc); err != nil {
		return err
	}
	return nil
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(err)
	}
}

func transformAndTransferRunes(
	reader io.Reader,
	writer io.Writer,
	inputBuffer *bytes.Buffer,
	outputBuffer *bytes.Buffer,
	transformFunc func(rune) rune,
) error {
	consecutiveErroneousInitialRune := false
	for _, err := inputBuffer.ReadFrom(reader); ; _, err = inputBuffer.ReadFrom(reader) {
		if errors.Is(err, io.EOF) {
			if consecutiveErroneousInitialRune {
				return ErrUnableToTransformRune
			}
			break
		}
		if err != nil {
			return err
		}

		err = transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)

		switch {
		case err == nil, errors.Is(err, ErrErroneousRune): // something was transformed, so write it
			consecutiveErroneousInitialRune = false
			if _, err = outputBuffer.WriteTo(writer); err != nil {
				return err
			}
			outputBuffer.Reset()
		case errors.Is(err, ErrErroneousInitialRune):
			if consecutiveErroneousInitialRune {
				return ErrUnableToTransformRune
			}
			consecutiveErroneousInitialRune = true
		default:
			return err
		}
	}
	return nil
}

var ErrUnableToTransformRune = errors.New("after two consecutive reads, could not transform rune")

// The input buffer is expected to be ready to be read from.
// The output buffer is expected to be ready to be written to.
// At the end, the input buffer is prepared to be written to again.
func transformRuneBuffers(inputBuffer *bytes.Buffer, outputBuffer *bytes.Buffer, transformFunc func(r rune) rune) error {
	iterCount := 0
	for inputRune, inputRuneSize, err := inputBuffer.ReadRune(); ; inputRune, inputRuneSize, err = inputBuffer.ReadRune() {
		if errors.Is(err, io.EOF) { // The whole input buffer has been read so end.
			inputBuffer.Reset()
			break
		}
		if err != nil {
			return err // Unexpected error has occurred.
		}

		if isInvalidRune(inputRune, inputRuneSize) {
			if err = inputBuffer.UnreadRune(); err != nil {
				return err // Unexpected error has occurred.
			}
			compactBuffer(inputBuffer)
			if iterCount == 0 {
				return ErrErroneousInitialRune
			}
			return ErrErroneousRune
		}

		transformedRune := transformFunc(inputRune)
		if _, err = outputBuffer.WriteRune(transformedRune); err != nil {
			return err
		}
		iterCount++
	}
	return nil
}

func isInvalidRune(r rune, size int) bool {
	return r == '\uFFFD' && size == 1
}

var ErrErroneousInitialRune = errors.New("invalid initial rune has been read")

func compactBuffer(inputBuffer *bytes.Buffer) {
	unreadChunk := inputBuffer.Bytes()
	inputBuffer.Reset()
	inputBuffer.Write(unreadChunk)
}

var ErrErroneousRune = errors.New("invalid rune has been read")
