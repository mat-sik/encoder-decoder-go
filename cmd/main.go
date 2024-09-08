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
	cipher, err := ciphers.NewCipher(argMap)
	if err != nil {
		panic(err)
	}
	ciphers.Run(cipher)
}
