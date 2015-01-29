package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
)

type Service struct {
	Timestamp string
	Subway    struct {
		Line []struct {
			Name   string
			Status string
			Text   string
			Date   string
			Time   string
		}
	}
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
	bodyBytes := make([]byte, 1000000) 
	n, err := resp.Body.Read(bodyBytes)
	if err != nil {
		log.Fatalf("Error reading response body: %i %v", n, err)
	}
	log.Printf("Read %v bytes", n)
	err = xml.Unmarshal(bodyBytes, &service)
	if err != nil {
		log.Fatalf("Error unmarshaling XML data: %v", err)
	}
	log.Println(service)
}
