package cipher

import (
	"bytes"
	"fmt"
	"github.com/mat-sik/encoder-decoder/internal/algorithms"
	"github.com/mat-sik/encoder-decoder/internal/parser"
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode/utf8"
)

func Test_newCipher(t *testing.T) {
	// given
	argMap := map[string]string{
		"-m": "encode",
		"-i": "foo.txt",
		"-o": "bar.txt",
		"-a": "mirror",
		"-k": "123",
	}
	expectedInput := &MirrorCipherInput{
		CipherInput: &CipherInput{
			Mode:    parser.Encode,
			Alg:     parser.Mirror,
			InPath:  "foo.txt",
			OutPath: "bar.txt",
		},
	}
	var expectedErr error = nil
	// when
	resultCipher, resultErr := newCipher(argMap)
	// then
	assert.Equal(t, expectedErr, resultErr)
	assert.IsType(t, expectedInput, resultCipher)
	assert.Equal(t, expectedInput, resultCipher)
}

func Test_transformRunes(t *testing.T) {
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
	err := transformRunes(inputBuffer, outputBuffer, transformFunc)
	inputBuffer.WriteByte(136)
	err = transformRunes(inputBuffer, outputBuffer, transformFunc)
	fmt.Println(err)
}
