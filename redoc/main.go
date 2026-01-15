// Package redoc ReDoc UI Middleware
//
// This package was copied and modified from https://github.com/go-openapi/runtime/blob/master/middleware/redoc.go
package redoc

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

// Opts configures the Redoc middleware
type Opts struct {
	BasePath string
	SpecURL  string
	SpecPath string
	RedocURL string
	Title    string
}

// EnsureDefaults in case some options are missing
func (r *Opts) EnsureDefaults() {
	if r.BasePath == "" {
		r.BasePath = "/"
	}
	if r.SpecURL == "" {
		r.SpecURL = "/swagger.json"
	}
	if r.SpecPath == "" {
		r.SpecPath = "swagger.json"
	}
	if r.RedocURL == "" {
		r.RedocURL = redocLatest
	}
	if r.Title == "" {
		r.Title = "API documentation"
	}
}

// sendSpec attempts to respond with the Swagger specification. If the specification file can't be found,
// it will respond with a 404.
func sendSpec(opts Opts, resp *echo.Response) {

	// Open the spec file.
	b, err := os.ReadFile(opts.SpecPath)
	if err != nil {
		resp.Header().Set("Content-Type", "text/plain")
		resp.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(resp, "%s not found", opts.SpecURL)
		return
	}

	// Send the contents of the spec file.
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	_, _ = resp.Write(b)
}

// Serve Serves Redoc documentation at either `/docs` or `/docs/`.
func Serve(opts Opts) echo.MiddlewareFunc {
	opts.EnsureDefaults()

	// Determine the valid documentation paths.
	docPath := path.Join(opts.BasePath, "docs")
	paths := map[string]string{
		docPath:       "",
		docPath + "/": "",
	}

	// Build the base HTML to return.
	tmpl := template.Must(template.New("redoc").Parse(redocTemplate))
	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, opts)
	b := buf.Bytes()

	// Return the middleware function.
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			resp := c.Response()

			// Return the documentation if the path is a documentation path.
			if _, ok := paths[req.URL.Path]; ok {
				resp.Header().Set("Content-Type", "text/html; charset=utf8")
				resp.WriteHeader(http.StatusOK)
				_, _ = resp.Write(b)
				return nil
			}

			// Return the spec if the path matches the spec URL path.
			if req.URL.Path == opts.SpecURL {
				sendSpec(opts, resp)
				return nil
			}

			// We didn't match the
			return next(c)
		}
	}
}

const (
	redocLatest   = "https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js"
	redocTemplate = `<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Title }}</title>
    <!-- needed for adaptive design -->
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!--
    ReDoc doesn't change outer page styles
    -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='{{ .SpecURL }}'></redoc>
    <script src="{{ .RedocURL }}"> </script>
  </body>
</html>
`
)
