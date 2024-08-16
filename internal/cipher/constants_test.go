package cipher 

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func Test_getModeValue(t *testing.T) {
    // given
    argMap := map[string]string {
        "-m": "encode",
    }
    expectedMode := Encode
    var expectedErr error = nil
    // when
    resultMode, resultErr := getModeValue(argMap)
    // then
    assert.Equal(t, expectedMode, resultMode)
    assert.Equal(t, expectedErr, resultErr)
}
