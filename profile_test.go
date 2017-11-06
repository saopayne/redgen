package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	expectedProfile := CreateDefaultProfile("")

	err := WriteProfileToFile(expectedProfile, defaultProfileName)
	if err != nil {
		t.Error(err.Error())
	}

	profileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, defaultProfileName))
	if err != nil {
		t.Error(err.Error())
	}

	profile := Profile{}
	err = json.Unmarshal(profileBytes, &profile)

	if !reflect.DeepEqual(expectedProfile, profile) {
		t.Errorf("Expected:\n %+v,\n\ngot:\n %+v", expectedProfile, profile)
	}
}


func TestSanitizeName(t *testing.T) {
	sampleName := " Sample Name "
	expectedResult := "sample_name"
	actualResult := SanitizeName(sampleName)
	assert.EqualValues(t, expectedResult, actualResult)
}
