package cipher

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
	expectedInput := &MirrorCipherInput{
		CipherInput: &CipherInput{
			Mode:    parser.Encode,
			Alg:     parser.Mirror,
			InPath:  "foo.txt",
			OutPath: "bar.txt",
		},
	}
	var expectedErr error = nil
	// when
	resultCipher, resultErr := newCipher(argMap)
	// then
	assert.Equal(t, expectedErr, resultErr)
	assert.IsType(t, expectedInput, resultCipher)
	assert.Equal(t, expectedInput, resultCipher)
}
