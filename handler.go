package main

import (
	"bytes"
	"io"
	"net/http"
	url2 "net/url"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
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
		m, ok := modules.Find(r.URL.Path)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if rest := strings.TrimPrefix(r.URL.Path, m.name()); strings.TrimSuffix(rest, "/") != "" {
			url, err := url2.ParseRequestURI(m.Readme)
			if err != nil {
				logrus.Errorf("parse readme url: %v", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			baseParts := strings.Split(url.Path, "/")
			url.Path = path.Join(append(baseParts[:len(baseParts)-1], rest)...)
			res, err := http.Get(url.String())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(res.StatusCode)
			copy(w, res)
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusSeeOther)
			return
		}
		w.Header().Set("content-type", "text/html")
		if err := moduleTemplate.Execute(w, m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func copy(w http.ResponseWriter, res *http.Response) {
	defer res.Body.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, res.Body)
	w.Write(buf.Bytes())
}
