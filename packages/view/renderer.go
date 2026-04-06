package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	nethttp "net/http"
	"os"
	"path/filepath"
	"strings"
)

// Renderer renders HTML templates from the demo application's resources/views directory.
type Renderer struct {
	basePath string
	funcs    template.FuncMap
}

// NewRenderer constructs a template renderer rooted at basePath.
func NewRenderer(basePath string) *Renderer {
	return &Renderer{
		basePath: filepath.Clean(basePath),
		funcs: template.FuncMap{
			"asset": func(assetPath string) string {
				if strings.HasPrefix(assetPath, "/") {
					return assetPath
				}

				return "/build/" + strings.TrimLeft(assetPath, "/")
			},
		},
	}
}

// Render renders a named view to a string.
func (r *Renderer) Render(name string, data any) (string, error) {
	target := filepath.Join(r.basePath, name+".html.tmpl")

	if _, err := os.Stat(target); err != nil {
		return "", fmt.Errorf("view: template %q not found", name)
	}

	files, err := r.templateFiles()
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(filepath.Base(target)).Funcs(r.funcs).ParseFiles(files...)
	if err != nil {
		return "", fmt.Errorf("view: parse templates: %w", err)
	}

	var output bytes.Buffer

	if err := tmpl.ExecuteTemplate(&output, filepath.Base(target), data); err != nil {
		return "", fmt.Errorf("view: render %q: %w", name, err)
	}

	return output.String(), nil
}

// WriteHTML writes a rendered view to the response.
func (r *Renderer) WriteHTML(writer nethttp.ResponseWriter, status int, name string, data any) error {
	rendered, err := r.Render(name, data)
	if err != nil {
		return err
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(rendered))

	return nil
}

func (r *Renderer) templateFiles() ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(r.basePath, func(file string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return nil
		}

		if strings.HasSuffix(entry.Name(), ".tmpl") {
			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("view: walk templates: %w", err)
	}

	return files, nil
}
