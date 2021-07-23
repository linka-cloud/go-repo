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

const pkg = `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="{{ .Import }} git https://{{ .Repository }}">
        <meta name="go-source" content="{{ .Import }} https://{{ .Repository }} https://{{ .Repository }}/tree/master{/dir} https://{{ .Repository }}/tree/master{/dir}/{file}#L{line}">
        {{/* <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/{{ .Import }}"> */}}
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" />
        <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.3/highlight.min.js" integrity="sha512-tHQeqtcNWlZtEh8As/4MmZ5qpy0wj04svWFK7MIzLmUVIzaHXS8eod9OmHxyBL1UET5Rchvw7Ih4ZDv5JojZww==" crossorigin="anonymous"></script>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.3/styles/monokai-sublime.min.css" integrity="sha512-8irWeigPA1Pm20tgaynUtqbAQ/zvOizMj7Olu0sF9kQTGabFfvlAHUqhslzHwr7OZO6Z0IN6VoXAALzipXIwgA==" crossorigin="anonymous" />
    </head>
    <body>
    <div class="container">
            <div class="row">
        {{/* Nothing to see here. Please <a href="https://pkg.go.dev/{{ .Import }}">move along</a>. */}}
        {{ .ReadmeHTML }}
            </div>
    </div>

    <script>hljs.initHighlightingOnLoad();</script>
    </body>
</html>
`
