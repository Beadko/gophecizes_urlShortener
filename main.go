package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshort "github.com/Beadko/gophecizes_urlShortener/internal"
)

func readFile(string) ([]byte, error) {
	f, err := os.ReadFile(*urlshort.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file: %w", err)
	}
	return f, nil
}

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	m, err := readFile(*urlshort.File)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	h, err := urlshort.Filehandler(m, mapHandler)
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
