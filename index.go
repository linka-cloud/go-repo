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
