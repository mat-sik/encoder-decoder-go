package ciphers

import (
	"github.com/mat-sik/encoder-decoder/internal/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newCipher(t *testing.T) {
	// given
	argMap := map[string]string{
		"-m": "encode",
		"-i": "foo.txt",
		"-o": "bar.txt",
		"-a": "mirror",
		"-k": "123",
	}
	expectedInput := &BasicCipherRunner{
		cipher: &MirrorCipherInput{
			CipherInput: &CipherInput{
				InPath:  "foo.txt",
				OutPath: "bar.txt",
			},
		},
		mode: parser.Encode,
	}
	var expectedErr error = nil
	// when
	resultCipher, resultErr := NewCipherRunner(argMap)
	// then
	assert.Equal(t, expectedErr, resultErr)
	assert.IsType(t, expectedInput, resultCipher)
	assert.Equal(t, expectedInput, resultCipher)
}
