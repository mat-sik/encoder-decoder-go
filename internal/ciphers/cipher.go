package ciphers

import (
	"bytes"

	"github.com/mat-sik/encoder-decoder/internal/algorithms"
	"github.com/mat-sik/encoder-decoder/internal/parser"
	"github.com/mat-sik/encoder-decoder/internal/transformer"
)

type CipherRunner interface {
	Run() error
}

type BasicCipherRunner struct {
	cipher Cipher
	mode   parser.Mode
}

func (cipherRunner *BasicCipherRunner) Run() error {
	cipher := cipherRunner.cipher
	switch cipherRunner.mode {
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
}

func NewCipherRunner(argMap map[string]string) (CipherRunner, error) {
	alg, err := parser.GetAlgValue(argMap)
	if err != nil {
		return nil, err
	}
	mode, err := parser.GetModeValue(argMap)
	if err != nil {
		return nil, err
	}
	var cipher Cipher
	switch alg {
	case parser.Caesar:
		cipher, err = newCaesarCipherInput(argMap)
	case parser.Mirror:
		cipher, err = newMirrorCipherInput(argMap)
	default:
		panic("technically this is not possible")
	}
	return &BasicCipherRunner{cipher, mode}, err
}

type CipherInput struct {
	InPath  string
	OutPath string
}

func newCipherInput(argMap map[string]string) (*CipherInput, error) {
	in, err := parser.GetInValue(argMap)
	if err != nil {
		return nil, err
	}
	out, err := parser.GetOutValue(argMap)
	if err != nil {
		return nil, err
	}
	return &CipherInput{in, out}, nil
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

func (input *CipherInput) transform(transformFunc func(rune) rune) error {
	inPath := input.InPath
	outPath := input.OutPath

	inBuffer := bytes.NewBuffer(make([]byte, 0, transformer.ReadBufferSize))
	outBuffer := bytes.NewBuffer(make([]byte, 0, transformer.WriteBufferSize))

	return transformer.FilesApplyFuncAndTransfer(inPath, outPath, inBuffer, outBuffer, transformFunc)
}
