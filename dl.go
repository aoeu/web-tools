package main

import (
	"path"
	"flag"
	"log"
	"net/url"
	"net/http"
	"os"
	"io"
)

func main() {
	flag.Parse()
	u := flag.Arg(0)
	_, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	name := path.Base(path.Clean(u))
	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(u)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
