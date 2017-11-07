package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	expectedProfile := CreateDefaultProfile("")

	err := WriteProfileToFile(expectedProfile, defaultProfilePath, defaultProfileName)
	if err != nil {
		t.Error(err.Error())
	}

	profileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, defaultProfileName))
	if err != nil {
		t.Error(err.Error())
	}

	profile := Profile{}
	err = json.Unmarshal(profileBytes, &profile)

	assert.EqualValues(t, expectedProfile, profile)
}

func TestSanitizeName(t *testing.T) {
	sampleName := " Sample Name "
	expectedResult := "sample_name"
	actualResult := SanitizeName(sampleName)
	assert.EqualValues(t, expectedResult, actualResult)
}
