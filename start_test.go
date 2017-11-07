package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmdInit(t *testing.T) {
	result, err := CmdInit()
	expected := "Default profile file created into ./profiles"
	assert.EqualValues(t, expected, result)
	assert.EqualValues(t, err, nil)
}

func TestCmdGenerate(t *testing.T) {
	actualResult, err := CmdGenerate("")
	assert.Equal(t, helpMsg, actualResult, "They should be equal")
	assert.EqualValues(t, nil, err)
	actualResult, err = CmdGenerate("file")
	expectedResult := "file file created into ./profiles"
	assert.EqualValues(t, expectedResult, actualResult, "They should be equal")
	assert.EqualValues(t, nil, err)
}
