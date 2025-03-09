package main

import (
	"flag"
	"fmt"
	"net/http"

	urlshort "github.com/Beadko/gophecizes_urlShortener/internal"
)

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	h, err := urlshort.FileHandler(mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")

	http.ListenAndServe(":8080", h)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
