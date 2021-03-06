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
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"go.uber.org/multierr"
)

var (
	mu sync.RWMutex
)

type Config struct {
	Modules []*Module
}

type Module struct {
	Import     string        `json:"import"`
	Repository string        `json:"repository"`
	Readme     string        `json:"readme"`
	ReadmeHTML template.HTML `json:"-"`
}

func (m Module) name() string {
	u, err := url.Parse("https://" + m.Import)
	if err != nil {
		return ""
	}
	return u.Path
}

func (m *Module) LoadReadme() error {
	if m.Readme == "" {
		m.Readme = m.Repository + "/raw/master/README.md"
	}
	if !strings.HasPrefix(m.Readme, "http") {
		m.Readme = "https://" + m.Readme
	}
	res, err := http.Get(m.Readme)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode > 300 {
		return fmt.Errorf("%s: request failed: %s (%d)", m.name(), res.Status, res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)
	m.ReadmeHTML = template.HTML(markdown.ToHTML(b, mdParser, nil))
	return nil
}

func NewModules(path string) (Modules, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c.Modules, nil
}

type Modules []*Module

func (m *Modules) Sort() {
	sort.Slice(*m, func(i, j int) bool {
		return strings.Compare((*m)[i].Import, (*m)[j].Import) > 0
	})
}

func (m *Modules) Find(name string) (*Module, bool) {
	for _, v := range *m {
		if strings.HasPrefix(name, v.name()) {
			return v, true
		}
	}
	return nil, false
}

func (m *Modules) LoadReadme() error {
	mu.Lock()
	defer mu.Unlock()
	errs := make(chan error, len(*m))
	for _, v := range *m {
		go func(mod *Module) {
			errs <- mod.LoadReadme()
		}(v)
	}
	var err error
	for range *m {
		err = multierr.Append(err, <-errs)
	}
	return err
}
