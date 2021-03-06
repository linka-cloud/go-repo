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
	_ "embed"
	"html/template"
)

var (
	//go:embed templates/index.gohtml
	index string
	//go:embed templates/package.gohtml
	pkg string
	//go:embed templates/header.gohtml
	header string
)

var (
	indexTemplate  = template.Must(template.New("index.gohtml").Parse(index))
	moduleTemplate = template.Must(template.New("package.gohtml").Parse(pkg))
	headerTemplate = template.Must(template.New("header.gohtml").Parse(header))
)
