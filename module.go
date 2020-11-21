package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
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
	Import     string `json:"import"`
	Repository string `json:"repository"`
	Readme     string `json:"readme"`
	ReadmeHTML template.HTML `json:"-"`
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
		return fmt.Errorf("request failed: %s (%d)", res.Status, res.StatusCode)
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

func (m *Modules) Find(name string) (*Module, bool) {
	for _, v := range *m {
		mod := path.Base(v.Import)
		if mod == name {
			return v, true
		}
	}
	return nil, false
}

func (m *Modules) LoadReadme() error {
	mu.RLock()
	defer mu.RUnlock()
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
