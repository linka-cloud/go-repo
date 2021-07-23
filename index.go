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

const index = `
<!DOCTYPE html>
<html>
    <head>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" />
    </head>
    <body>
        <div class="container">
            <div class="row">
                <table class="u-full-width">
                    <thead>
                        <tr>
                            <th>Package</th>
                            <th>Source</th>
                            <th>Documentation</th>
                        </tr>
                    </thead>
                    <tbody>

                        {{- range . }}
                        <tr>
                            <td>{{ .Import }}</td>
                            <td>
                                <a href="//{{ .Repository }}">{{ .Repository }}</a>
                            </td>
                            <td>
                                <a href="//pkg.go.dev/{{ .Import }}">
                                    <img src="//img.shields.io/badge/godoc-reference-blue?style=for-the-badge" alt="GoDoc" />
                                </a>
                            </td>
                        </tr>
                        {{- end }}
                    </tbody>
                </table>
            </div>
        </div>
    </body>
</html>
`
