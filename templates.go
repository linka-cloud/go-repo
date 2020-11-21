package main

import (
	"html/template"
)

var (
	indexTemplate  = template.Must(template.New("index.gohtml").Parse(index))
	moduleTemplate = template.Must(template.New("package.gohtml").Parse(pkg))
)
