package input

type Alg string

const (
	Caesar Alg = "caesar"
	Mirror Alg = "mirror"
)

func NewAlg(algString string) (Alg, error) {
	switch Alg(algString) {
	case Caesar:
		return Caesar, nil
	case Mirror:
		return Mirror, nil
	default:
		return "", &ErrUnknownAlgorithm{algString}
	}
}

type ErrUnknownAlgorithm struct {
	Alg string
}

func (e *ErrUnknownAlgorithm) Error() string {
	return "unknown algorithm: " + e.Alg
}

type Mode string

const (
	Encode Mode = "encode"
	Decode Mode = "decode"
)

func NewMode(modeString string) (Mode, error) {
	switch Mode(modeString) {
	case Encode:
		return Encode, nil
	case Decode:
		return Decode, nil
	default:
		return "", &ErrUnknownMode{modeString}
	}
}

type ErrUnknownMode struct {
	Mode string
}

func (e *ErrUnknownMode) Error() string {
	return "unknown mode: " + e.Mode
}

type Flag string

const (
	ChosenMode     Flag = "-m"
	ChosenModeFull Flag = "--mode"
	In             Flag = "-i"
	InFull         Flag = "--input"
	Out            Flag = "-o"
	OutFull        Flag = "--output"
	ChosenAlg      Flag = "-a"
	ChosenAlgFull  Flag = "--algorithm"
	Key            Flag = "-k"
	KeyFull        Flag = "--key"
)
