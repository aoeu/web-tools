package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// TODO: Is there an idiom for naming variable comprised entirely of acronyms?
var APIURL = "https://www.googleapis.com/urlshortener/v1/url"
var userProvidedURL = flag.String("url", "", "A URL to shorten.")

func shortenURL(longURL string) (shortURL string, err error) {
	_, err = url.Parse(longURL)
	if err != nil {
		return
	}
	payload := strings.NewReader(`{"longUrl": "` + longURL + `"}`)
	resp, err := http.Post(APIURL, "application/json", payload)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	type Message struct {
		Kind, Id, LongUrl string
	}
	dec := json.NewDecoder(resp.Body)
	var m Message
	err = dec.Decode(&m)
	shortURL = m.Id
	return
}

func main() {
	flag.Parse()
	shortURL, err := shortenURL(*userProvidedURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(shortURL)
}

// curl https://www.googleapis.com/urlshortener/v1/url -H 'Content-Type: application/json' -d '{"longUrl": "http://www.google.com/"}'

/*
{
 "kind": "urlshortener#url",
 "id": "http://goo.gl/fbsS",
 "longUrl": "http://www.google.com/"
}
*/
