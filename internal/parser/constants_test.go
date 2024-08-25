package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getModeValue(t *testing.T) {
	// given
	argMap := map[string]string{
		"-m": "encode",
	}
	expectedMode := Encode
	var expectedErr error = nil
	// when
	resultMode, resultErr := GetModeValue(argMap)
	// then
	assert.Equal(t, expectedMode, resultMode)
	assert.Equal(t, expectedErr, resultErr)
}
