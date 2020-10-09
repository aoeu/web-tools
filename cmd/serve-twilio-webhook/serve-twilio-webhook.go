// serve-twilio-webhook is an HTTP web server (that also supports FastCGI)
// that functions as a webhook for inbound SMS from Twilio.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/fcgi"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"
)

var (
	logFile *os.File
	nyc     *time.Location
	re      = regexp.MustCompile("Body\t(.*)")
)

func inscribe(s string) {
	if logFile != nil {
		ss := fmt.Sprintf("%v : %v\n", time.Now().In(nyc), s)
		if _, err := logFile.WriteString(ss); err != nil {
			log.Fatalf("could not write string '%v' to logfile", s)
		}
	} else {
		log.Println(s)
	}
}

func HandleHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respond(err, http.StatusInternalServerError, rw)
		}
		s, err := format(b)
		if err != nil {
			respond(err, http.StatusInternalServerError, rw)
		}
		inscribe(s)
	case "GET", "PUT", "DELETE":
		rw.WriteHeader(http.StatusNotImplemented)
	}
}

func format(twilioQueryParams []byte) (formatted string, err error) {
	s, err := url.QueryUnescape(string(twilioQueryParams))
	if err != nil {
		s := "could not unescape query params '%v' : %v"
		return "", fmt.Errorf(s, string(twilioQueryParams), err)
	}
	s = strings.Replace(s, "=", "\t", -1)
	eightSpaces := "        "
	s = strings.Replace(s, "&", "\n"+eightSpaces, -1)
	s = eightSpaces + s
	msg := re.FindStringSubmatch(s)
	if len(msg) > 1 {
		formatted += msg[1] + "\n"
	}
	buf := bytes.NewBufferString(formatted)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, s)
	w.Flush()
	return buf.String(), nil
}

// webServer is a type that exists only because the fcgi package does not have
// a function that models the http.HandleFunc function (and instead only
// has a function that mirrors http.ListenAndServe function).
type webServer struct{}

func (s webServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	HandleHTTP(rw, r)
}

func respond(err error, statusCode uint, rw http.ResponseWriter) {
	log.Printf("could not read request body : %v\n", err)
	rw.WriteHeader(http.StatusInternalServerError)
}

func main() {
	args := struct {
		useFCGI     bool
		port        string
		logFilepath string
	}{}

	flag.BoolVar(&args.useFCGI, "fcgi", false, "serve via  FCGI (instead of HTTP)")
	flag.StringVar(&args.port, "port", ":8080", "the port to serve on (exclusively for HTTP, not FCGI)")
	flag.StringVar(&args.logFilepath, "log", os.Args[0]+".txt", "the path of a file to log inbound payloads to")
	flag.Parse()

	p, err := filepath.Abs(args.logFilepath)
	if err != nil {
		s := "could not parse log file path '%v' : %v"
		log.Fatalf(s, args.logFilepath, err)
	}
	logFile, err = os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		s := "could not create log file at path '%v' : %v"
		log.Fatalf(s, p, err)
	}
	tz := "America/New_York"
	nyc, err = time.LoadLocation(tz)
	if err != nil {
		s := "could not load time.Location for timezone named '%v' : %v"
		log.Fatalf(s, tz, err)
	}

	if args.useFCGI {
		if err := fcgi.Serve(nil, webServer{}); err != nil {
			inscribe(err.Error())
		}
	} else {
		if err := http.ListenAndServe(args.port, webServer{}); err != nil {
			inscribe(err.Error())
		}
	}
}