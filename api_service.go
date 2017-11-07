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

type LibrarianService struct {
	sling    *sling.Sling
	URL      string
	Auth     string
	HttpType string
}

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

func SendHTTPRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
