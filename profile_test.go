package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	expectedProfile := CreateDefaultProfile()

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

func TestLoadProfile(t *testing.T) {
	expectedProfile := CreateDefaultProfile()

	profileBytes := bytes.NewBufferString("{\"name\":\"DefaultProfile\",\"baseDailyConsumption\":100,\"hourlyProfiles\":" +
		"{\"00\":1,\"01\":1,\"02\":1,\"03\":1,\"04\":1,\"05\":1,\"06\":1,\"07\":1,\"08\":1,\"09\":1,\"10\":1,\"11\":1,\"12\":1," +
		"\"13\":1,\"14\":1,\"15\":1,\"16\":1,\"17\":1,\"18\":1,\"19\":1,\"20\":1,\"21\":1,\"22\":1,\"23\":1},\"WeeklyProfiles\":" +
		"{\"Fri\":1,\"Mon\":1,\"Sat\":1,\"Sun\":1,\"Thu\":1,\"Tue\":1,\"Wed\":1},\"monthlyProfiles\":{\"Apr\":1,\"Aug\":1,\"Dec\":1," +
		"\"Feb\":1,\"Jan\":1,\"Jul\":1,\"Jun\":1,\"Mar\":1,\"May\":1,\"Nov\":1,\"Oct\":1,\"Sep\":1},\"variability\":5,\"unit\":\"kW\"," +
		"\"interval\":15,\"startAt\":{\"year\":2017,\"month\":\"Jan\",\"hour\":6}, \"readings\":[]}")

	profile, err := NewProfileFromJson(profileBytes.Bytes())
	if err != nil {
		t.Error(err.Error())
	}

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
