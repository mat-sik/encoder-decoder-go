package parser

import (
	"fmt"
	"strings"
)

func Parse(args []string) (map[string]string, error) {
	argMap := make(map[string]string)
	for position, arg := range args {
		if err := parseArg(&argMap, position, arg); err != nil {
			return nil, err
		}
	}
	return argMap, nil
}

func parseArg(argMap *map[string]string, position int, arg string) error {
	if !isValidArg(arg) {
		return &ErrInvalidArg{arg, position}
	}
	flag, value, isPairArg := parsePairArg(arg)
	if !isPairArg {
		flag = arg
	}
	(*argMap)[flag] = value
	return nil
}

func isValidArg(arg string) bool {
	return strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--")
}

func parsePairArg(arg string) (string, string, bool) {
	pair := strings.SplitN(arg, "=", 2)
	if len(pair) != 2 {
		return "", "", false
	}
	return pair[0], pair[1], true
}

type ErrInvalidArg struct {
	arg      string
	position int
}

func (e *ErrInvalidArg) Error() string {
	return fmt.Sprintf("invalid argument: %s at: %d", e.arg, e.position)
}
