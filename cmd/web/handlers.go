package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"new-portfolio/internal/models"
)

func (app *application) sectionIndex(section, page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caser := cases.Title(language.English)
		data := app.NewPageTemplateData(r)
		data.Title = caser.String(section)
		data.Description = fmt.Sprintf("Latest from %s.", section)
		data.Data = models.ListPageData{
			Section: section,
			Entries: app.store.List(section),
		}
		app.render(w, r, page, data)
	}
}

func (app *application) sectionShow(section, page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			app.notFound(w)
			return
		}

		showData, err := app.store.FindForShow(section, slug)
		if err != nil {
			app.notFound(w)
			return
		}

		data := app.NewPageTemplateData(r)
		data.Title = showData.Entry.Title
		data.Description = showData.Entry.Summary
		data.Data = showData

		app.render(w, r, page, data)
	}
}

func (app *application) sectionAddComment(section string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimSpace(r.PathValue("slug"))
		if slug == "" {
			app.notFound(w)
			return
		}

		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		name := r.PostForm.Get("name")
		content := r.PostForm.Get("content")

		var parentID *int64
		parentIDRaw := strings.TrimSpace(r.PostForm.Get("parent_id"))
		if parentIDRaw != "" {
			parsed, err := strconv.ParseInt(parentIDRaw, 10, 64)
			if err != nil || parsed <= 0 {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			parentID = &parsed
		}

		if err := app.store.AddComment(section, slug, name, content, parentID); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/%s#comments", section, slug), http.StatusSeeOther)
	}
}

func (app *application) sectionArchive(section, page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caser := cases.Title(language.English)
		data := app.NewPageTemplateData(r)
		data.Title = fmt.Sprintf("%s archive", caser.String(section))
		data.Description = fmt.Sprintf("Archive for %s.", section)
		data.Data = models.ArchivePageData{
			Section: section,
			Entries: app.store.List(section),
		}
		app.render(w, r, page, data)
	}
}

func (app *application) sectionArchiveYear(section, page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year, err := parseYear(r.PathValue("year"))
		if err != nil {
			app.notFound(w)
			return
		}

		caser := cases.Title(language.English)
		entries := app.store.ListByYear(section, year)
		data := app.NewPageTemplateData(r)
		data.Title = fmt.Sprintf("%s archive %d", caser.String(section), year)
		data.Description = fmt.Sprintf("%s archive for %d.", section, year)
		data.Data = models.ArchivePageData{
			Section:     section,
			ArchiveYear: year,
			Entries:     entries,
		}
		app.render(w, r, page, data)
	}
}

func (app *application) sectionTag(section, page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := strings.TrimSpace(r.PathValue("tag"))
		if tag == "" {
			app.notFound(w)
			return
		}

		caser := cases.Title(language.English)
		data := app.NewPageTemplateData(r)
		data.Title = fmt.Sprintf("%s tagged %s", caser.String(section), tag)
		data.Description = fmt.Sprintf("Entries in %s tagged %s.", section, tag)
		data.Data = models.ArchivePageData{
			Section: section,
			Tag:     tag,
			Entries: app.store.ListByTag(section, tag),
		}
		app.render(w, r, page, data)
	}
}

func (app *application) setTheme(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	theme := r.Form.Get("theme")
	if theme != "auto" && theme != "light" && theme != "dark" {
		theme = "auto"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    theme,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((365 * 24 * time.Hour).Seconds()),
	})

	back := r.Referer()
	if back == "" {
		back = "/"
	}
	http.Redirect(w, r, back, http.StatusSeeOther)
}

func (app *application) atom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?><feed xmlns="http://www.w3.org/2005/Atom"><title>Feed placeholder</title></feed>`))
}

func (app *application) sitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://localhost:4000/</loc></url></urlset>`))
}

func (app *application) robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("User-agent: *\nAllow: /\n"))
}
