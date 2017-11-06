package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestProfileIsEqual(t *testing.T) {
	expectedProfile := CreateDefaultProfile("")

	profileBytes := bytes.NewBufferString("{\"name\":\"DefaultProfile\",\"baseDailyConsumption\":18," +
		"\"hourlyProfiles\":{\"00\":1,\"01\":1,\"02\":1,\"03\":1,\"04\":1,\"05\":1,\"06\":1,\"07\":1,\"08\":1," +
		"\"09\":1,\"10\":1,\"11\":1,\"12\":1,\"13\":1,\"14\":1,\"15\":1,\"16\":1,\"17\":1,\"18\":1,\"19\":1,\"20\":1," +
		"\"21\":1,\"22\":1,\"23\":1},\"WeeklyProfiles\":{\"Fri\":1,\"Mon\":1,\"Sat\":1,\"Sun\":1,\"Thu\":1,\"Tue\":1," +
		"\"Wed\":1},\"monthlyProfiles\":{\"Apr\":1,\"Aug\":1,\"Dec\":1,\"Feb\":1,\"Jan\":1,\"Jul\":1,\"Jun\":1,\"Mar\":1," +
		"\"May\":1,\"Nov\":1,\"Oct\":1,\"Sep\":1},\"variability\":5,\"unit\":\"kW\",\"interval\":15,\"startAt\":" +
		"{\"year\":2017,\"month\":\"Jan\",\"day\":1,\"hour\":6}, \"readings\":[]}")

	profile, err := NewProfileFromJson(profileBytes.Bytes())
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(expectedProfile, profile) {
		t.Errorf("Expected:\n %+v,\n\ngot:\n %+v", expectedProfile, profile)
	}
}

func TestIsUnitValid(t *testing.T) {
	// a non-existent value is passed
	unitToCheck := "YW"
	actualValidity := IsUnitValid(unitToCheck)
	assert.EqualValues(t, true, actualValidity)

	// pass in an existent value
	unitToCheck = "mW"
	actualValidity = IsUnitValid(unitToCheck)
	assert.True(t, actualValidity)
}

func TestIsValueInList(t *testing.T) {
	namesList := []string{"ademola", "oyewale", "valid"}

	stringToFind := "ademol"
	valueExists := IsValueInList(stringToFind, namesList)
	// name doesn't exist in list
	assert.EqualValues(t, false, valueExists)

	// name exists
	stringToFind = "valid"
	valueExists = IsValueInList(stringToFind, namesList)
	assert.EqualValues(t, true, valueExists)
}

func TestValidateProfile(t *testing.T) {
	profileBytesValid := bytes.NewBufferString("{\"name\":\"DefaultProfile\",\"baseDailyConsumption\":100," +
		"\"hourlyProfiles\":{\"00\":1,\"01\":1,\"02\":1,\"03\":1,\"04\":1,\"05\":1,\"06\":1,\"07\":1,\"08\":1,\"09\":1," +
		"\"10\":1,\"11\":1,\"12\":1,\"13\":1,\"14\":1,\"15\":1,\"16\":1,\"17\":1,\"18\":1,\"19\":1,\"20\":1,\"21\":1," +
		"\"22\":1,\"23\":1},\"WeeklyProfiles\":{\"Fri\":1,\"Mon\":1,\"Sat\":1,\"Sun\":1,\"Thu\":1,\"Tue\":1,\"Wed\":1}," +
		"\"monthlyProfiles\":{\"Apr\":1,\"Aug\":1,\"Dec\":1,\"Feb\":1,\"Jan\":1,\"Jul\":1,\"Jun\":1,\"Mar\":1,\"May\":1," +
		"\"Nov\":1,\"Oct\":1,\"Sep\":1},\"variability\":5,\"unit\":\"kW\",\"interval\":15,\"startAt\":{\"year\":2000, " +
		"\"month\":\"Jan\",\"hour\":6}, \"readings\":[]}")
	profile, err := NewProfileFromJson(profileBytesValid.Bytes())
	if err != nil {
		t.Error(err.Error())
	}
	err = profile.ValidateProfile()
	assert.Empty(t, err)
}

func TestValidateName(t *testing.T) {
	profileBytesWithInvalidName := bytes.NewBufferString("{\"name\":\"De\",\"baseDailyConsumption\":100,\"hourlyProfiles\":" +
		"{\"00\":1,\"01\":1,\"02\":1,\"03\":1,\"04\":1,\"05\":1,\"06\":1,\"07\":1,\"08\":1,\"09\":1,\"10\":1,\"11\":1,\"12\":1," +
		"\"13\":1,\"14\":1,\"15\":1,\"16\":1,\"17\":1,\"18\":1,\"19\":1,\"20\":1,\"21\":1,\"22\":1,\"23\":1},\"WeeklyProfiles\":" +
		"{\"Fri\":1,\"Mon\":1,\"Sat\":1,\"Sun\":1,\"Thu\":1,\"Tue\":1,\"Wed\":1},\"monthlyProfiles\":{\"Apr\":1,\"Aug\":1,\"Dec\":1," +
		"\"Feb\":1,\"Jan\":1,\"Jul\":1,\"Jun\":1,\"Mar\":1,\"May\":1,\"Nov\":1,\"Oct\":1,\"Sep\":1},\"variability\":5,\"unit\":\"kW\"," +
		"\"interval\":15,\"startAt\":{\"year\":2000, \"month\":\"Jan\",\"hour\":6}, \"readings\":[]}")
	profileBytesWithValidName := bytes.NewBufferString("{\"name\":\"Default\",\"baseDailyConsumption\":100,\"hourlyProfiles\":" +
		"{\"00\":-1,\"01\":1,\"02\":1,\"03\":1,\"04\":1,\"05\":1,\"06\":1,\"07\":1,\"08\":1,\"09\":1,\"10\":1,\"11\":1,\"12\":1," +
		"\"13\":1,\"14\":1,\"15\":1,\"16\":1,\"17\":1,\"18\":1,\"19\":1,\"20\":1,\"21\":1,\"22\":1,\"23\":1},\"WeeklyProfiles\"" +
		":{\"Fri\":1,\"Mon\":1,\"Sat\":1,\"Sun\":1,\"Thu\":1,\"Tue\":1,\"Wed\":1},\"monthlyProfiles\":{\"Apr\":1,\"Aug\":1,\"Dec\":" +
		"1,\"Feb\":1,\"Jan\":1,\"Jul\":1,\"Jun\":1,\"Mar\":1,\"May\":1,\"Nov\":1,\"Oct\":1,\"Sep\":1},\"variability\":5,\"unit\":" +
		"\"kW\",\"interval\":15,\"startAt\":{\"year\":2000, \"month\":\"Jan\",\"hour\":6}, \"readings\":[]}")
	// invalid name
	profile, err := NewProfileFromJson(profileBytesWithInvalidName.Bytes())
	if err != nil {
		t.Error(err.Error())
	}
	err = ValidateName(profile)
	if err != nil {
		assert.Error(t, err, "")
	}
	// Valid name
	profile, err = NewProfileFromJson(profileBytesWithValidName.Bytes())
	if err != nil {
		t.Error(err.Error())
	}
	err = ValidateName(profile)
	assert.Empty(t, err)
}

func TestValidateHourlyProfiles(t *testing.T) {
	profileBytesWithInvalidHour := bytes.NewBufferString("{\"name\":\"Default\",\"baseDailyConsumption\":100," +
		"\"hourlyProfiles\":{\"00\":-1,\"01\":1,\"02\":1,\"03\":1,\"04\":1,\"05\":1,\"06\":1,\"07\":1,\"08\":1,\"09\"" +
		":1,\"10\":1,\"11\":1,\"12\":1,\"13\":1,\"14\":1,\"15\":1,\"16\":1,\"17\":1,\"18\":1,\"19\":1,\"20\":1,\"21\":1," +
		"\"22\":1,\"23\":1},\"WeeklyProfiles\":{\"Fri\":1,\"Mon\":1,\"Sat\":1,\"Sun\":1,\"Thu\":1,\"Tue\":1,\"Wed\":1}," +
		"\"monthlyProfiles\":{\"Apr\":1,\"Aug\":1,\"Dec\":1,\"Feb\":1,\"Jan\":1,\"Jul\":1,\"Jun\":1,\"Mar\":1,\"May\":1," +
		"\"Nov\":1,\"Oct\":1,\"Sep\":1},\"variability\":5,\"unit\":\"kW\",\"interval\":15,\"startAt\":{\"year\":2000, " +
		"\"month\":\"Jan\",\"hour\":6}, \"readings\":[]}")
	profileBytesWithValidHour := bytes.NewBufferString("{\"name\":\"Default\",\"baseDailyConsumption\":100," +
		"\"hourlyProfiles\":{\"00\":2,\"01\":1,\"02\":1,\"03\":1,\"04\":1,\"05\":1,\"06\":1,\"07\":1,\"08\":1," +
		"\"09\":1,\"10\":1,\"11\":1,\"12\":1,\"13\":1,\"14\":1,\"15\":1,\"16\":1,\"17\":1,\"18\":1,\"19\":1,\"20\":1," +
		"\"21\":1,\"22\":1,\"23\":1},\"WeeklyProfiles\":{\"Fri\":1,\"Mon\":1,\"Sat\":1,\"Sun\":1,\"Thu\":1,\"Tue\":1," +
		"\"Wed\":1},\"monthlyProfiles\":{\"Apr\":1,\"Aug\":1,\"Dec\":1,\"Feb\":1,\"Jan\":1,\"Jul\":1,\"Jun\":1,\"Mar\":1," +
		"\"May\":1,\"Nov\":1,\"Oct\":1,\"Sep\":1},\"variability\":5,\"unit\":\"kW\",\"interval\":15,\"startAt\":{\"year\":2000," +
		" \"month\":\"Jan\",\"hour\":6}, \"readings\":[]}")
	// invalid hour
	profile, err := NewProfileFromJson(profileBytesWithInvalidHour.Bytes())
	if err != nil {
		t.Error(err.Error())
	}
	err = ValidateHourlyProfiles(profile)
	if err != nil {
		assert.Error(t, err, "")
	}
	// valid hour
	profile, err = NewProfileFromJson(profileBytesWithValidHour.Bytes())
	if err != nil {
		t.Error(err.Error())
	}
	err = ValidateHourlyProfiles(profile)
	assert.Empty(t, err)
}
