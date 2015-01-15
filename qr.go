package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(QR))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

type QRParams struct {
	Url, Dimensions string
}

func QR(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w,
		QRParams{
			Url:        req.FormValue("url"),
			Dimensions: req.FormValue("dimensions"),
		})
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
<br>
<br>
<form action="/" name=f method="GET"><input maxLength=1024 size=70 name=url value="" title="Text to QR Encode"><input type=submit value="Show QR" name=qr>
<select name=dimensions>
<option value="100x100">Small</option>
<option value="300x300" selected="selected">Medium</option>
<option value="500x500">Large</option>
</select>
</form>
{{if .Url}}
<a href="{{.Url}}">{{.Url}}</a>
<br>
<img src="http://chart.apis.google.com/chart?chs={{.Dimensions}}&cht=qr&choe=UTF-8&chl={{.Url}}" />
{{end}}
</body>
</html>
`
