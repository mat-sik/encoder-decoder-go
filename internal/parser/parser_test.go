package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	// given
	input := []string{
		"-arg",
		"-one=two",
		"--three=four",
		"--five",
	}
	var expectedErr error = nil
	expectedMap := map[string]string{
		"-arg":    "",
		"-one":    "two",
		"--three": "four",
		"--five":  "",
	}
	// when
	resultMap, resultErr := parse(input)
	// then
	assert.Equal(t, expectedErr, resultErr)
	assert.Equal(t, expectedMap, resultMap)
}

func TestParse_NoArgs(t *testing.T) {
	// given
	var input []string
	var expectedErr error = nil
	expectedMap := make(map[string]string)
	// when
	resultMap, resultErr := parse(input)
	// then
	assert.Equal(t, expectedErr, resultErr)
	assert.Equal(t, expectedMap, resultMap)
}

func TestParse_IncorrectArg(t *testing.T) {
	// given
	input := []string{
		"-arg",
		"-one=two",
		"wrong",
		"--five",
	}
	var expectedErr error = &ErrInvalidArg{"wrong", 2}
	var expectedMap map[string]string = nil
	// when
	resultMap, resultErr := parse(input)
	// then
	assert.Equal(t, expectedErr, resultErr)
	assert.Equal(t, expectedMap, resultMap)
}

func TestIsProperArg_NoHyphen(t *testing.T) {
	// given
	input := "arg"
	expected := false
	// when
	result := isValidArg(input)
	// then
	assert.Equal(t, expected, result)
}

func TestIsProperArg_SingleHyphen(t *testing.T) {
	// given
	input := "-arg"
	expected := true
	// when
	result := isValidArg(input)
	// then
	assert.Equal(t, expected, result)
}

func TestIsProperArg_DoubleHyphen(t *testing.T) {
	// given
	input := "--hello"
	expected := true
	// when
	result := isValidArg(input)
	// then
	assert.Equal(t, expected, result)
}

func TestIsPairArg_GoodInput(t *testing.T) {
	// given
	input := "-arg=one"
	expectedKey := "-arg"
	expectedValue := "one"
	expectedIsPair := true
	// when
	resultKey, resultValue, resultIsPair := parsePairArg(input)
	// then
	assert.Equal(t, expectedKey, resultKey)
	assert.Equal(t, expectedValue, resultValue)
	assert.Equal(t, expectedIsPair, resultIsPair)
}

func TestIsPairArg_NoPairArg(t *testing.T) {
	// given
	input := "-argone"
	expectedKey := ""
	expectedValue := ""
	expectedIsPair := false
	// when
	resultKey, resultValue, resultIsPair := parsePairArg(input)
	// then
	assert.Equal(t, expectedKey, resultKey)
	assert.Equal(t, expectedValue, resultValue)
	assert.Equal(t, expectedIsPair, resultIsPair)
}
