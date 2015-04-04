package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

// "Echo" the HTTP headers and payload of the HTTP request in the response.
func main() {
	port := flag.String("port", ":8080", "The HTTP port to listen on.")
	url := flag.String("url", "/", "The URL path to serve.")
	flag.Parse()
	http.HandleFunc(*url, func(w http.ResponseWriter, r *http.Request) {
		m := io.MultiWriter(os.Stdout, w)
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = m.Write(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
