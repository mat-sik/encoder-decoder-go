package cipher

import "github.com/mat-sik/encoder-decoder/internal/parser"

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
