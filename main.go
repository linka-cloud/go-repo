package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type Module struct {
	Package    string
	Repository string
	Readme     template.HTML
}

func (m *Module) LoadReadme() error {
	url := path.Join(m.Repository, "/raw/master/README.md")
	res, err := http.Get("https://" + url)
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
	m.Readme = template.HTML(markdown.ToHTML(b, mdParser, nil))
	return nil
}

type Modules []*Module

func (m *Modules) Find(pkg string) (*Module, bool) {
	for _, v := range *m {
		if v.Package == pkg {
			return v, true
		}
	}
	return nil, false
}

func (m *Modules) LoadReadme() error {
	errs := make(chan error, len(*m))
	for _, v := range *m {
		go func(mod *Module) {
			errs <- mod.LoadReadme()
		}(v)
	}
	var err error
	for range *m {
		multierr.Append(err, <-errs)
	}
	return err
}

var (
	modules = Modules{
		{
			Package:    "go.adphi.net/gonextcloud",
			Repository: "git.adphi.net/adphi/gonextcloud",
		},
		{
			Package:    "go.adphi.net/gogs-cli",
			Repository: "git.adphi.net/adphi/gogs-cli",
		},
	}
	indexTemplate  = template.Must(template.New("index.gohtml").ParseFiles("./index.gohtml"))
	moduleTemplate = template.Must(template.New("package.gohtml").ParseFiles("./package.gohtml"))
)

func main() {
	go func() {
		if err := modules.LoadReadme(); err != nil {
			logrus.WithError(err).Error("failed to load all readme")
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("content-type", "text/html")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		switch r.URL.Path {
		case "/":
			if err := indexTemplate.Execute(w, modules); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			m, ok := modules.Find("go.adphi.net" + r.URL.Path)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("content-type", "text/html")
			if err := moduleTemplate.Execute(w, m); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
	if err := http.ListenAndServe(":8888", nil); err != nil {
		logrus.Fatal(err)
	}
}
