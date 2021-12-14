// Copyright 2021 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		if r.Form.Get("go-get") == "1" {
			w.Header().Set("content-type", "text/html")
			if err := headerTemplate.Execute(w, m); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
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
