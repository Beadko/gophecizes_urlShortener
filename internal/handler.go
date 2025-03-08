package urlshort

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	File = flag.String("file", "paths.yaml", "Path to url file")
)

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if redirectURL, exists := pathsToUrls[path]; exists {
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// FileHandler will parse the provided file based on type and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the file, then the
// fallback http.Handler will be called instead.
func FileHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathURLs []pathURL
	pathsToUrls := make(map[string]string)

	ext := filepath.Ext(*File)
	switch ext {
	case ".yaml":
		if err := yaml.Unmarshal(data, &pathURLs); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, &pathURLs); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unknown file type %s", ext)
	}

	for _, p := range pathURLs {
		pathsToUrls[p.Path] = p.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}
