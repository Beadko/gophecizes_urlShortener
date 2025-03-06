package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshort "github.com/Beadko/gophecizes_urlShortener/internal"
)

var (
	file = flag.String("file", "paths.yaml", "Path to url file")
)

func readFile(string) ([]byte, error) {
	f, err := os.ReadFile(*file)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file: %w", err)
	}
	return f, nil
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yaml, err := readFile(*file)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
