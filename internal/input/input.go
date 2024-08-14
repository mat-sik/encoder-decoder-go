package input

import (
	"strconv"
)

type CipherInput struct {
	Mode    Mode
	Alg     Alg
	InPath  string
	OutPath string
}

type Cipher interface {
	encode()
	decode()
}

type CaesarCipherInput struct {
	CipherInput     CipherInput
	CaesarCipherKey int
}

func (input *CaesarCipherInput) encode() {}
func (input *CaesarCipherInput) decode() {}

type MirrorCipherInput CipherInput

func (input *MirrorCipherInput) encode() {}
func (input *MirrorCipherInput) decode() {}

func getAlg(argMap map[string]string) (Alg, error) {
	algString, ok := argMap[string(ChosenAlg)]
	if !ok {
		algString, ok = argMap[string(ChosenAlgFull)]
	}
	if !ok {
		return "", &ErrUnspecifiedMode{}
	}
	alg, err := NewAlg(algString)
	if err != nil {
		return "", err
	}
	return alg, nil
}

type ErrUnspecifiedMode struct {
}

func (e *ErrUnspecifiedMode) Error() string {
	return "Unspecified Mode"
}

func NewCaesarCipherInput(argMap map[string]string) (Cipher, error) {
	alg, err := getAlg(argMap)
	if err != nil {
		return nil, err
	}
	switch alg {
	case Caesar:
		return nil, nil
	case Mirror:
		return nil, nil
	default:
		return nil, nil
	}
}

func NewMirrorCipherInput(argMap map[string]string) (Cipher, error) {
	input := CipherInput{}
	for key, arg := range argMap {
		if err := input.setValueFromFlag(key, arg); err != nil {
			return nil, err
		}
	}
	// is it copy?
	mirrorCipherInput := MirrorCipherInput(input)
	return &mirrorCipherInput, nil
}

func (input *CipherInput) setValueFromFlag(key string, value string) error {
	switch Flag(key) {
	case ChosenMode, ChosenModeFull:
		mode, err := NewMode(value)
		if err != nil {
			return err
		}
		input.Mode = mode
	case ChosenAlg, ChosenAlgFull:
		alg, err := NewAlg(value)
		if err != nil {
			return err
		}
		input.Alg = alg
	case In, InFull:
		input.InPath = value
	case Out, OutFull:
		input.OutPath = value
	default:
		return &ErrUnsupportedFlag{key}
	}
	return nil
}

func (input *CaesarCipherInput) setValueFromFlag(key string, value string) error {
	switch Flag(key) {
	case Key, KeyFull:
		cipherKey, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		input.CaesarCipherKey = cipherKey
	default:
		return &ErrUnsupportedFlag{key}
	}
	return nil
}

type ErrUnsupportedFlag struct {
	flag string
}

func (e ErrUnsupportedFlag) Error() string {
	return "unsupported flag: " + e.flag
}
