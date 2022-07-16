package handler

import (
	"encoding/json"
	"net/http"

	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nextUrl, present := pathsToUrls[r.RequestURI]
		if present {
			http.Redirect(w, r, nextUrl, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	entries, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(entries)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYAML(yml []byte) ([]entry, error) {
	var entries []entry
	err := yaml.Unmarshal(yml, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func buildMap(entries []entry) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, e := range entries {
		pathsToUrls[e.Path] = e.Url
	}
	return pathsToUrls
}

type entry struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//		[{"path": "/json","url": "https://jsoneditoronline.org/"}]
//
// The only errors that can be returned all related to having
// invalid JSO data.
//
func JSONHandler(js []byte, fallback http.Handler) (http.HandlerFunc, error) {
	entries, err := parseJson(js)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(entries)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseJson(js []byte) ([]entry, error) {
	var entries []entry
	if err := json.Unmarshal(js, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// BoldDBHandler will open the BoltDB pointed by the provided path
// and return an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding URL.
//  If the path doesn't exist in the DB, then the fallback http.Handler
// will be called instead.
//
// DB is expected to have as paths as key and the url to redirect as value
func BoltDBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("routes"))
			if nextUrl := b.Get([]byte(r.RequestURI)); nextUrl != nil {
				http.Redirect(w, r, string(nextUrl), http.StatusMovedPermanently)
			} else {
				fallback.ServeHTTP(w, r)
			}
			return nil
		})
	}
}
