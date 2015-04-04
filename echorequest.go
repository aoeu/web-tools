package main

import (
	"flag"
	"fmt"
	"io"
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
		log.Println("Received request:")
		for key, value := range r.Header {
			fmt.Fprintf(m, "%v : %v\n", key, value)
		}
		_, err := io.Copy(m, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
