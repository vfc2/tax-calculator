package main

import (
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
)

type Views struct {
	templates map[string]*template.Template
}

func NewViews() (*Views, error) {
	tpl, err := loadTemplates()
	if err != nil {
		return nil, err
	}

	return &Views{
		templates: tpl,
	}, nil
}

func (v Views) render(w http.ResponseWriter, template string, name string, data any, logger *slog.Logger) {
	tpl := v.templates[template]
	if tpl == nil {
		logger.Error("the requested template could not be found", "template", template)
		return
	}

	err := tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		logger.Error("the requested template could not be executed", "template", template, "error", err)
	}
}

func loadTemplates() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	tpl, err := template.ParseFiles(
		"./assets/templates/layout.html",
		"./assets/templates/partials/tax_input.html",
	)
	if err != nil {
		return nil, err
	}

	cache["home"] = tpl

	err = filepath.WalkDir("./assets/templates/partials/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		tplName := d.Name()[:len(d.Name())-len(filepath.Ext(d.Name()))]

		tpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		cache[tplName] = tpl

		return nil
	})
	if err != nil {
		return nil, err
	}

	return cache, nil
}
