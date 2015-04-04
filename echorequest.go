package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


// "Echo" the HTTP headers and payload of the HTTP request in the response.
func main() {
	port := flag.String("port", ":8080", "The HTTP port to listen on.")
	url := flag.String("url", "/", "The URL path to serve.")
	flag.Parse()
	http.HandleFunc(*url, func(w http.ResponseWriter, r *http.Request) {
		m := io.MultiWriter(os.Stdout, w)
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println("Received request:")
		for key, value := range r.Header {
			fmt.Fprintf(m, "%v : %v\n", key, value)
		}
		fmt.Fprintf(m, string(b))
	})
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
