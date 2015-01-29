package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
)

type Service struct {
	Timestamp string `xml:"timestamp"`
	Subway    struct {
		Line []struct {
			Name   string `xml:"name"`
			Status string `xml:"status"`
			Text   string `xml:"text"`
			Date   string `xml:"date"`
			Time   string `xml:"time"`
		} `xml:"line"`
	} `xml:"subway"`
}

const rawURL = "http://mta.info/status/serviceStatus.txt"

func main() {
	MTAURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}
	resp, err := http.Get(MTAURL.String())
	if err != nil {
		log.Fatalf("Error getting URL: %v", err)
	}
	defer resp.Body.Close()
	var service Service
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	bodyBytes := make([]byte, 0)
	for _, rune := range body {
		if "U+0000" != fmt.Sprintf("%U", rune) {
			bodyBytes = append(bodyBytes, rune)
		}
	}

	err = xml.Unmarshal(bodyBytes, &service)
	if err != nil {
		log.Fatalf("Error unmarshaling XML data: %v", err)
	}
	log.Println(service)
}
