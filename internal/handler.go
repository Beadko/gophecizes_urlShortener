package urlshort

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v4/stdlib"
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
func FileHandler(fallback http.Handler) (http.HandlerFunc, error) {
	var pathURLs []pathURL
	pathsToUrls := make(map[string]string)

	ext := filepath.Ext(*File)
	switch ext {
	case ".yaml":
		m, err := readFile(*File)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		if err := yaml.Unmarshal(m, &pathURLs); err != nil {
			return nil, err
		}
	case ".json":
		m, err := readFile(*File)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		if err := json.Unmarshal(m, &pathURLs); err != nil {
			return nil, err
		}
	case ".db":
		db, err := sql.Open("pgx", "postgresql://root@localhost:26258/defaultdb?sslmode=disable")
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		defer db.Close()
		rows, err := db.QueryContext(context.Background(), "SELECT path, url FROM paths")
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var p pathURL
			if err := rows.Scan(&p.Path, &p.URL); err != nil {
				return nil, fmt.Errorf("%w", err)
			}
			pathURLs = append(pathURLs, p)
		}

	default:
		return nil, fmt.Errorf("Unknown file type %s", ext)
	}

	for _, p := range pathURLs {
		pathsToUrls[p.Path] = p.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func readFile(file string) ([]byte, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file: %w", err)
	}
	return f, nil
}
