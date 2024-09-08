package main

import (
	"os"

	"github.com/mat-sik/encoder-decoder/internal/ciphers"
	"github.com/mat-sik/encoder-decoder/internal/parser"
)

func main() {
	args := os.Args

	argMap, err := parser.Parse(args)
	if err != nil {
		panic(err)
	}
	cipherRunner, err := ciphers.NewCipherRunner(argMap)
	if err != nil {
		panic(err)
	}
	if err = cipherRunner.Run(); err != nil {
		panic(err)
	}
}
