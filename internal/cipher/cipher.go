package cipher

import (
	"bytes"
	"github.com/mat-sik/encoder-decoder/internal/parser"
	"io"
)

type Cipher interface {
	encode()
	decode()
}

func newCipher(argMap map[string]string) (Cipher, error) {
	alg, err := parser.GetAlgValue(argMap)
	if err != nil {
		return nil, err
	}
	switch alg {
	case parser.Caesar:
		return newCaesarCipherInput(argMap)
	case parser.Mirror:
		return newMirrorCipherInput(argMap)
	default:
		panic("technically this is not possible")
	}
}

type CipherInput struct {
	Mode    parser.Mode
	Alg     parser.Alg
	InPath  string
	OutPath string
}

func newCipherInput(argMap map[string]string) (*CipherInput, error) {
	mode, err := parser.GetModeValue(argMap)
	if err != nil {
		return nil, err
	}
	alg, err := parser.GetAlgValue(argMap)
	if err != nil {
		return nil, err
	}
	in, err := parser.GetInValue(argMap)
	if err != nil {
		return nil, err
	}
	out, err := parser.GetOutValue(argMap)
	if err != nil {
		return nil, err
	}
	return &CipherInput{mode, alg, in, out}, nil
}

type CaesarCipherInput struct {
	CipherInput     *CipherInput
	CaesarCipherKey int
}

func newCaesarCipherInput(argMap map[string]string) (*CaesarCipherInput, error) {
	cipherInput, err := newCipherInput(argMap)
	if err != nil {
		return nil, err
	}
	key, err := parser.GetIntKeyValue(argMap)
	if err != nil {
		return nil, err
	}
	return &CaesarCipherInput{cipherInput, key}, nil
}

func (input *CaesarCipherInput) encode() {}
func (input *CaesarCipherInput) decode() {}

type MirrorCipherInput struct {
	CipherInput *CipherInput
}

func newMirrorCipherInput(argMap map[string]string) (*MirrorCipherInput, error) {
	cipherInput, err := newCipherInput(argMap)
	if err != nil {
		return nil, err
	}
	return &MirrorCipherInput{cipherInput}, nil
}

func (input *MirrorCipherInput) encode() {}
func (input *MirrorCipherInput) decode() {}

const ReadBufferSize = 4 * 1024
const WriteBufferSize = 4 * ReadBufferSize

// The input buffer is expected to be ready to be read from.
// The output buffer is expected to be ready to be written to.
// At the end, the input buffer is prepared to be written to again.
func transformRunes(inputBuffer *bytes.Buffer, outputBuffer *bytes.Buffer, transformFunc func(r rune) rune) error {
	for inputRune, inputRuneSize, err := inputBuffer.ReadRune(); ; inputRune, inputRuneSize, err = inputBuffer.ReadRune() {
		if err != nil {
			if err == io.EOF { // The whole input buffer has been read so end.
				inputBuffer.Reset()
				break
			} else { // Unexpected error has occurred.
				panic(err)
			}
		}

		// Invalid or not whole rune has been read, so leave if for next transformRunes operation and end.
		if inputRuneSize == 1 && inputRune == '\uFFFD' {
			if err = inputBuffer.UnreadRune(); err != nil {
				panic(err)
			}
			unreadChunk := inputBuffer.Bytes()
			inputBuffer.Reset()
			inputBuffer.Write(unreadChunk)
			return &ErrErroneousRune{}
		} else { // Transform rune and write it to output buffer.
			transformedRune := transformFunc(inputRune)
			outputBuffer.WriteRune(transformedRune)
		}
	}
	return nil
}

type ErrErroneousRune struct{}

func (err *ErrErroneousRune) Error() string {
	return "invalid rune has been read"
}
