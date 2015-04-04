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

type multiWriter struct {
	writers []io.Writer
}

func newMultiWriter(writers ...io.Writer) *multiWriter {
	m := &multiWriter{writers: make([]io.Writer, len(writers))}
	for i, w := range writers {
		m.writers[i] = w
	}
	return m
}

func (m *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range m.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

var port string
var url string

// "Echo" the HTTP headers and payload of the HTTP request in the response.
func main() {
	flag.StringVar(&port, "port", ":8080", "The HTTP port to listen on.")
	flag.StringVar(&url, "url", "/", "The URL path to serve.")
	flag.Parse()
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		m := newMultiWriter(os.Stdout, w)
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
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
