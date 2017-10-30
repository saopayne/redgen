package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
)

const (
	UUID = "99999999-0000-1111-2222-999999999999"
	ApiUrl = "https://api.energyanalytics.eu:41115/readings/" + UUID
	Authorization = "Bearer: eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJqaXQiOiJmYjVhZGI0YjU5ODcwMWI4Iiwic2NwIjoibGEiLCJ2ZXIiOjEsInN1YiI6Ijk5OTk5OTk5LTAwMDAtMTExMS0yMjIyLTk5OTk5OTk5OTk5OSIsImV4cCI6MTUwOTk1MzA0Nn0.95WVy0oM-KqsTwt5TNnE87KtE4jDgPajgFEgbDvrzFWoiow__8T4mdDxetG5vbjwTSl1FPWq-Smea0tvT5dg"
	HeaderContentType = "application/vnd.api+json"
)

// ProfileService provides methods for creating readings.
type ProfileService struct {
	sling *sling.Sling
}

// Response is a simplified data the enectiva API sends
type Response struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type PostData struct {
	Data  Reading  `json:"data"`
}

type ProfileError struct {
	Message string `json:"message"`
	Errors  []struct {
		Resource string `json:"resource"`
		Field    string `json:"field"`
		Code     string `json:"code"`
	} `json:"errors"`
}

func (e ProfileError) Error() string {
	return fmt.Sprintf("enectiva reading: %v %+v", e.Message, e.Errors)
}

// This function creates the schedule that sends the readings to the API
// It accepts a function parameter which specifies the function to perform
// It also accepts a time delay to run the function in subsequent times (typically 15 minutes)
func sendReadingsSchedule(action func(), delay time.Duration) chan bool {
	stop := make(chan bool)
	ticker := time.NewTicker(delay)
	go func() {
		for {
			select {
			case <-ticker.C:
				action()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}

// Given a series, this function goes through the readings in the series
// For each given given, send the reading to the API
// If the POST was successful, remove from temp storage and avoid retry
func (s *ProfileService) sendReadingsAction(p Profile) (*http.Response, error) {
	for reading := range p.Readings {
		//readingError := new(ProfileError)
		eachReading := p.Readings[reading]
		eachReading.MeterId = "test"
		eachReading.Sender = "ademola"
		postData := &PostData{ Data: eachReading}
		jsonReading, err := json.Marshal(postData)
		fmt.Println(jsonReading)
		req, err := BuildHTTPRequest("POST", Authorization, ApiUrl, jsonReading)
		if err != nil {
			return nil, err
		}

		resp, err := SendHTTPRequest(req)
		if err != nil {
			return nil, err
		}
		fmt.Println("response Status:", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		return resp, err
	}
	return nil, nil
}

//BuildHTTPRequest builds a request (Sets Body and Header)
func BuildHTTPRequest(httpType string, auth string, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(httpType, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", HeaderContentType)
	req.Header.Add("Content-Type", HeaderContentType)
	req.Header.Add("Authorization", auth)
	return req, nil
}

//SendHTTPRequest sends the request and returns a response from SMC Server
func SendHTTPRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
