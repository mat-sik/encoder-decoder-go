package ciphers

import (
	"bytes"

	"github.com/mat-sik/encoder-decoder/internal/algorithms"
	"github.com/mat-sik/encoder-decoder/internal/parser"
	"github.com/mat-sik/encoder-decoder/internal/transformer"
)

func Run(cipher Cipher) error {
	switch cipher.getMode() {
	case parser.Encode:
		return cipher.encode()
	case parser.Decode:
		return cipher.decode()
	default:
		panic("technically this is not possible")
	}
}

type Cipher interface {
	encode() error
	decode() error
	getMode() parser.Mode
}

func NewCipher(argMap map[string]string) (Cipher, error) {
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
	InPath  string
	OutPath string
}

func newCipherInput(argMap map[string]string) (*CipherInput, error) {
	mode, err := parser.GetModeValue(argMap)
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
	return &CipherInput{mode, in, out}, nil
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

func (input *CaesarCipherInput) encode() error {
	var key = int32(input.CaesarCipherKey)
	encodeFunc := algorithms.NewOffsetRuneFunc(key)
	return input.CipherInput.transform(encodeFunc)
}

func (input *CaesarCipherInput) decode() error {
	var key = -int32(input.CaesarCipherKey)
	decodeFunc := algorithms.NewOffsetRuneFunc(key)
	return input.CipherInput.transform(decodeFunc)
}

func (input *CaesarCipherInput) getMode() parser.Mode {
	return input.CipherInput.Mode
}

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

func (input *MirrorCipherInput) encode() error {
	encodeFunc := algorithms.GetMirrorRuneLatin1
	return input.CipherInput.transform(encodeFunc)
}

func (input *MirrorCipherInput) decode() error {
	decodeFunc := algorithms.GetMirrorRuneLatin1
	return input.CipherInput.transform(decodeFunc)
}

func (input *MirrorCipherInput) getMode() parser.Mode {
	return input.CipherInput.Mode
}

func (input *CipherInput) transform(transformFunc func(rune) rune) error {
	inPath := input.InPath
	outPath := input.OutPath

	inBuffer := bytes.NewBuffer(make([]byte, 0, transformer.ReadBufferSize))
	outBuffer := bytes.NewBuffer(make([]byte, 0, transformer.WriteBufferSize))

	return transformer.FilesApplyFuncAndTransfer(inPath, outPath, inBuffer, outBuffer, transformFunc)
}
