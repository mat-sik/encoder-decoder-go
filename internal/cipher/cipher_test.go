package cipher

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_newCipher(t *testing.T) {
    // given
    argMap := map[string]string {
        "-m": "encode",
        "-i": "foo.txt",
        "-o": "bar.txt",
        "-a": "mirror",
        "-k": "123",
    }
    expectedInput := &MirrorCipherInput{
        CipherInput: &CipherInput{
            Mode: Encode,
            Alg: Mirror,
            InPath: "foo.txt",
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
