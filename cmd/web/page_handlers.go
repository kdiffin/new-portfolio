package main

import (
	"net/http"

	"new-portfolio/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.NewPageTemplateData(r)
	data.Title = "Home"
	data.Description = "A minimal publishing site built with Go and Tailwind CSS."
	data.Data = models.HomePageData{
		WritingEntries: app.store.List("writing"),
		MicroEntries:   app.store.List("micro"),
		Books:          app.store.LatestBooks(8),
	}
	app.render(w, r, "home.tmpl", data)
}

func (app *application) booksAndPapers(w http.ResponseWriter, r *http.Request) {
	data := app.NewPageTemplateData(r)
	data.Title = "Books and Papers"
	data.Description = "Books and papers reading log."
	data.Data = models.BooksPageData{
		Books: app.store.Books(),
	}
	app.render(w, r, "books_index.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.NewPageTemplateData(r)
	data.Title = "About"
	data.Description = "About this site."
	app.render(w, r, "about.tmpl", data)
}
