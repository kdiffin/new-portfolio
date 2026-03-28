package render

import (
	"fmt"
	"html/template"
	"path/filepath"
)

func NewTemplateCache(root string, funcs template.FuncMap) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(root, "pages", "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("glob pages: %w", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(funcs).ParseFiles(filepath.Join(root, "base.tmpl"))
		if err != nil {
			return nil, fmt.Errorf("parse base for %s: %w", name, err)
		}

		ts, err = ts.ParseGlob(filepath.Join(root, "partials", "*.tmpl"))
		if err != nil {
			return nil, fmt.Errorf("parse partials for %s: %w", name, err)
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("parse page %s: %w", name, err)
		}

		cache[name] = ts
	}

	return cache, nil
}
