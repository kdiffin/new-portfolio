package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /about", app.about)

	mux.HandleFunc("GET /writing", app.sectionIndex("writing", "writing_index.tmpl"))
	mux.HandleFunc("GET /writing/{slug}", app.sectionShow("writing", "writing_show.tmpl"))
	mux.HandleFunc("POST /writing/{slug}/comments", app.sectionAddComment("writing"))
	mux.HandleFunc("GET /writing/archive", app.sectionArchive("writing", "archive_index.tmpl"))
	mux.HandleFunc("GET /writing/archive/{year}", app.sectionArchiveYear("writing", "archive_year.tmpl"))
	mux.HandleFunc("GET /writing/tag/{tag}", app.sectionTag("writing", "tag_archive.tmpl"))

	mux.HandleFunc("GET /books", app.booksAndPapers)
	mux.HandleFunc("GET /books/{slug}", app.sectionShow("books", "books_show.tmpl"))
	mux.HandleFunc("POST /books/{slug}/comments", app.sectionAddComment("books"))
	mux.HandleFunc("GET /books/archive", app.sectionArchive("books", "archive_index.tmpl"))
	mux.HandleFunc("GET /books/archive/{year}", app.sectionArchiveYear("books", "archive_year.tmpl"))

	mux.HandleFunc("GET /projects", app.sectionIndex("projects", "projects_index.tmpl"))
	mux.HandleFunc("GET /projects/{slug}", app.sectionShow("projects", "projects_show.tmpl"))
	mux.HandleFunc("POST /projects/{slug}/comments", app.sectionAddComment("projects"))
	mux.HandleFunc("GET /projects/archive", app.sectionArchive("projects", "archive_index.tmpl"))

	mux.HandleFunc("GET /micro", app.sectionIndex("micro", "micro_index.tmpl"))
	mux.HandleFunc("GET /micro/{slug}", app.sectionShow("micro", "micro_show.tmpl"))
	mux.HandleFunc("POST /micro/{slug}/comments", app.sectionAddComment("micro"))
	mux.HandleFunc("GET /micro/archive", app.sectionArchive("micro", "archive_index.tmpl"))
	mux.HandleFunc("GET /micro/archive/{year}", app.sectionArchiveYear("micro", "archive_year.tmpl"))

	mux.HandleFunc("POST /theme", app.setTheme)

	mux.HandleFunc("GET /atom.xml", app.atom)
	mux.HandleFunc("GET /sitemap.xml", app.sitemap)
	mux.HandleFunc("GET /robots.txt", app.robots)

	return app.recoverPanic(app.logRequest(app.secureHeaders(mux)))
}
