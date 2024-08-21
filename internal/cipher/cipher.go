package cipher

type Cipher interface {
	encode()
	decode()
}

func newCipher(argMap map[string]string) (Cipher, error) {
	alg, err := getAlgValue(argMap)
	if err != nil {
		return nil, err
	}
	switch alg {
	case Caesar:
		return newCaesarCipherInput(argMap)
	case Mirror:
		return newMirrorCipherInput(argMap)
	default:
		panic("technically this is not possible")
	}
}

type CipherInput struct {
	Mode    Mode
	Alg     Alg
	InPath  string
	OutPath string
}

func newCipherInput(argMap map[string]string) (*CipherInput, error) {
	mode, err := getModeValue(argMap)
	if err != nil {
		return nil, err
	}
	alg, err := getAlgValue(argMap)
	if err != nil {
		return nil, err
	}
	in, err := getInValue(argMap)
	if err != nil {
		return nil, err
	}
	out, err := getOutValue(argMap)
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
	key, err := getIntKeyValue(argMap)
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
