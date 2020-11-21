package main

import (
	"net/http"
	"path"
)

func modulesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	switch r.URL.Path {
	case "/":
		if err := indexTemplate.Execute(w, modules); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		m, ok := modules.Find(path.Base(r.URL.Path))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("content-type", "text/html")
		if err := moduleTemplate.Execute(w, m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
