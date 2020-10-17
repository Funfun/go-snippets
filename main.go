package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", "0.0.0.0:8080", "host:port of listen http server")
)

func main() {
	flag.Parse()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
