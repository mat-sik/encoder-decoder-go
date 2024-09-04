package transformer

import (
	"bytes"
	"errors"
	"io"
	"os"
)

// WriteBufferSize by making it 4 times, we guarantee that it will be able to fit 4kb of transformed runes.
const (
	ReadBufferSize  = 4 * 1024
	WriteBufferSize = 4 * ReadBufferSize
)

func filesApplyFuncAndTransfer(
	inputFilePath string,
	outputFilePath string,
	inputBuffer *bytes.Buffer,
	outputBuffer *bytes.Buffer,
	transformFunc func(r rune) rune,
) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer safeCloseFile(inputFile)

	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer safeCloseFile(outputFile)

	return applyFuncAndTransfer(inputFile, outputFile, inputBuffer, outputBuffer, transformFunc)
}

func safeCloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(err)
	}
}

func applyFuncAndTransfer(
	reader io.Reader,
	writer io.Writer,
	inputBuffer *bytes.Buffer,
	outputBuffer *bytes.Buffer,
	transformFunc func(rune) rune,
) error {
	inputBufferCapacity := int64(inputBuffer.Cap())
	limitedReader := io.LimitedReader{R: reader, N: inputBufferCapacity}

	consecutiveInvalidRune := false
	for {
		copiedLimitedReader := limitedReader
		readSize, err := inputBuffer.ReadFrom(&copiedLimitedReader)
		if err != nil {
			return err
		}
		if readSize == 0 {
			if consecutiveInvalidRune || inputBuffer.Len() != 0 {
				return ErrUnableToTransformRune
			}
			break
		}

		err = runeBuffersApplyFuncAndTransfer(inputBuffer, outputBuffer, transformFunc)

		switch {
		case err == nil, errors.Is(err, ErrErroneousRune): // something was transformed, so write it
			consecutiveInvalidRune = false
			if _, err = outputBuffer.WriteTo(writer); err != nil {
				return err
			}
		case errors.Is(err, ErrErroneousInitialRune):
			if consecutiveInvalidRune {
				return ErrUnableToTransformRune
			}
			consecutiveInvalidRune = true
		default:
			return err
		}
	}
	return nil
}

// The input buffer is expected to be ready to be read from.
// The output buffer is expected to be ready to be written to.
// At the end, the input buffer is prepared to be written to again.
func runeBuffersApplyFuncAndTransfer(
	inputBuffer *bytes.Buffer,
	outputBuffer *bytes.Buffer,
	transformFunc func(r rune) rune,
) error {
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

func compactBuffer(inputBuffer *bytes.Buffer) {
	unreadChunk := inputBuffer.Bytes()
	inputBuffer.Reset()
	inputBuffer.Write(unreadChunk)
}

var (
	ErrUnableToTransformRune = errors.New("after two consecutive reads, could not transform rune")
	ErrErroneousInitialRune  = errors.New("invalid initial rune has been read")
	ErrErroneousRune         = errors.New("invalid rune has been read")
)
