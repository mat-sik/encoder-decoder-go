package main

import (
	"os"

	"github.com/mat-sik/encoder-decoder/internal/cipher"
	"github.com/mat-sik/encoder-decoder/internal/parser"
)

func main() {
    args := os.Args

    argMap, err := parser.Parse(args)
    if err != nil {
        panic(err)
    }
    c, err := cipher.NewCipher(argMap)
    if err != nil {
        panic(err)
    }
    cipher.Run(c)
}
