package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dghubble/sling"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	UUID              = "99999999-0000-1111-2222-999999999999"
	ApiUrl            = "https://api.energyanalytics.eu:41115/readings/" + UUID
	Authorization     = "Bearer: eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJqaXQiOiJmYjVhZGI0YjU5ODcwMWI4Iiwic2NwIjoibGEiLCJ2ZXIiOjEsInN1YiI6Ijk5OTk5OTk5LTAwMDAtMTExMS0yMjIyLTk5OTk5OTk5OTk5OSIsImV4cCI6MTUwOTk1MzA0Nn0.95WVy0oM-KqsTwt5TNnE87KtE4jDgPajgFEgbDvrzFWoiow__8T4mdDxetG5vbjwTSl1FPWq-Smea0tvT5dg"
	HeaderContentType = "application/vnd.api+json"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// LibrarianService provides methods for sending readings to the API
type LibrarianService struct {
	sling    *sling.Sling
	URL      string
	Auth     string
	HttpType string
}

// Response is a simplified data the enectiva API sends
type Response struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type LibrarianError struct {
	Message string `json:"message"`
	Errors  []struct {
		Resource string `json:"resource"`
		Field    string `json:"field"`
		Code     string `json:"code"`
	} `json:"errors"`
}

func (e LibrarianError) Error() string {
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
func (s *LibrarianService) sendReadingsAction(p Profile) (*http.Response, error) {
	service := &LibrarianService{
		URL:      ApiUrl,
		Auth:     Authorization,
		HttpType: http.MethodPost,
	}
	totalReadings := len(p.Readings)

	eachReading := p.Readings[totalReadings-1]
	eachReading.MeterId = "test"
	eachReading.Sender = "ademola"
	postData := map[string]Reading{"data": eachReading}
	jsonReading, err := json.Marshal(postData)
	req, err := service.BuildHTTPRequest(jsonReading)
	if err != nil {
		return nil, err
	}
	resp, err := SendHTTPRequest(httpClient, req)
	if err != nil {
		return nil, err
	}
	fmt.Println("response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("response Body:", string(body))
		return resp, err
	}

	return resp, nil
}

//BuildHTTPRequest builds a request (Sets Body and Header)
func (s *LibrarianService) BuildHTTPRequest(body []byte) (*http.Request, error) {
	req, err := http.NewRequest(s.HttpType, s.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", HeaderContentType)
	req.Header.Add("Content-Type", HeaderContentType)
	req.Header.Add("Authorization", s.Auth)
	return req, nil
}

//SendHTTPRequest sends the request and returns a response from SMC Server
func SendHTTPRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
