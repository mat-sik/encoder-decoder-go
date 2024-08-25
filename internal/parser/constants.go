package parser

import (
	"fmt"
	"strconv"
)

type Alg string

const (
	Caesar Alg = "caesar"
	Mirror Alg = "mirror"
)

func newAlg(algString string) (Alg, error) {
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

func newMode(modeString string) (Mode, error) {
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

func GetIntKeyValue(argMap map[string]string) (int, error) {
	intString, err := getKeyValue(argMap)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(intString)
}

func GetInValue(argMap map[string]string) (string, error) {
	return getFlagValue(argMap, In, InFull)
}

func GetOutValue(argMap map[string]string) (string, error) {
	return getFlagValue(argMap, Out, OutFull)
}

func getKeyValue(argMap map[string]string) (string, error) {
	return getFlagValue(argMap, Key, KeyFull)
}

func GetModeValue(argMap map[string]string) (Mode, error) {
	return getMappedValue(argMap, ChosenMode, ChosenModeFull, newMode)
}

func GetAlgValue(argMap map[string]string) (Alg, error) {
	return getMappedValue(argMap, ChosenAlg, ChosenAlgFull, newAlg)
}

func getMappedValue[T ~string](argMap map[string]string, flag Flag, fullFlag Flag, newConst func(string) (T, error)) (T, error) {
	value, err := getFlagValue(argMap, flag, fullFlag)
	if err != nil {
		return "", err
	}
	mapped, err := newConst(value)
	if err != nil {
		return mapped, err
	}
	return mapped, nil
}

func getFlagValue(argMap map[string]string, flag Flag, fullFlag Flag) (string, error) {
	flagString := string(flag)
	value, ok := argMap[flagString]
	if !ok {
		flagString = string(fullFlag)
		value, ok = argMap[flagString]
	}
	if !ok {
		return value, &ErrMissingFlag{flag, fullFlag}
	}
	return value, nil
}

type ErrMissingFlag struct {
	RequiredFlag     Flag
	RequiredFlagFull Flag
}

func (err *ErrMissingFlag) Error() string {
	return fmt.Sprintf("required flag: %s or %s is missing", err.RequiredFlag, err.RequiredFlagFull)
}
