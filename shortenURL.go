package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"encoding/json"
)

// TODO: Is there an idiom for naming variable comprised entirely of acronyms?
var APIURL = "https://www.googleapis.com/urlshortener/v1/url"
var longURL = flag.String("url", "", "A URL to shorten.")

func main() {
	flag.Parse()
	payload := strings.NewReader(`{"longUrl": "` + *longURL + `"}`)
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
	if err := dec.Decode(&m); err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.Id)
}

// curl https://www.googleapis.com/urlshortener/v1/url -H 'Content-Type: application/json' -d '{"longUrl": "http://www.google.com/"}'

/*
{
 "kind": "urlshortener#url",
 "id": "http://goo.gl/fbsS",
 "longUrl": "http://www.google.com/"
}
*/

