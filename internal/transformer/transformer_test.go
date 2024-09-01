package transformer

import (
	"bytes"
	"fmt"
	"github.com/mat-sik/encoder-decoder/internal/algorithms"
	"testing"
	"unicode/utf8"
)

func Test_transformRuneBuffers(t *testing.T) {
	// given
	airplane := '\u2708'
	byteAirplane := make([]byte, utf8.RuneLen(airplane))
	utf8.EncodeRune(byteAirplane, airplane)

	inputBuffer := new(bytes.Buffer)
	inputBuffer.WriteRune(airplane)
	inputBuffer.Write(byteAirplane[:2])

	outputBuffer := new(bytes.Buffer)

	transformFunc := func(r rune) rune {
		return algorithms.OffsetRuneForward(r, 10)
	}
	// when
	err := transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)
	inputBuffer.WriteByte(136)
	err = transformRuneBuffers(inputBuffer, outputBuffer, transformFunc)
	fmt.Println(err)
}
