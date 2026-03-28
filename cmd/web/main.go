package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"new-portfolio/internal/content"
	"new-portfolio/internal/models"
	"new-portfolio/internal/render"
)

type application struct {
	logger        *slog.Logger
	store         *content.Store
	templateCache map[string]*template.Template
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	templateFuncs := template.FuncMap{
		"fmtDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"title": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"themeClass": themeClass,
		"limit": func(entries []*models.Entry, max int) []*models.Entry {
			if max <= 0 || len(entries) <= max {
				return entries
			}
			return entries[:max]
		},
		"isActive": func(currentPath, href string) bool {
			if href == "/" {
				return currentPath == "/"
			}
			return currentPath == href || strings.HasPrefix(currentPath, href+"/")
		},
		"safeHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
	}

	templateCache, err := render.NewTemplateCache(filepath.Join("ui", "html"), templateFuncs)
	if err != nil {
		logger.Error("could not build template cache", slog.Any("error", err))
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		store:         content.NewStore(),
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("starting server", slog.String("addr", srv.Addr))
	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}
}

func (app *application) NewTemplateData(r *http.Request) *models.TemplateData {
	data := &models.TemplateData{
		Path:     r.URL.Path,
		Year:     time.Now().Year(),
		Sections: app.store.Sections(),
	}

	if cookie, err := r.Cookie("theme"); err == nil {
		data.Theme = cookie.Value
	} else {
		data.Theme = "auto"
	}

	return data
}

func (app *application) render(w http.ResponseWriter, r *http.Request, page string, data *models.TemplateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s not found", page))
		return
	}

	if data == nil {
		data = app.NewTemplateData(r)
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = buf.WriteTo(w)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logger.Error("server error", slog.Any("error", err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func themeClass(theme string) string {
	switch theme {
	case "light":
		return "light"
	default:
		return "dark"
	}
}

func parseYear(value string) (int, error) {
	y, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if y < 1900 || y > 2200 {
		return 0, errors.New("year out of range")
	}
	return y, nil
}
