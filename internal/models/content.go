package models

import (
	"html/template"
	db "new-portfolio/internal/db/sqlc"
	"time"
)

// Non-page models
type Entry struct {
	Slug        string
	Title       string
	Summary     string
	BodyHTML    template.HTML
	Section     string
	PublishedAt time.Time
	Tags        []string
	Comments    []db.GetCommentsBySlugRow
}

type BookReview struct {
	Entry
	Author     string
	FinishedAt time.Time
	Rating     int
}

type SectionLink struct {
	Label string
	Href  string
}

// Page models
type PageTemplateData struct {
	Title       string // for seo
	Description string // for seo
	Path        string // for the navbar highlighting
	Theme       string
	Year        int           // for the footer
	Sections    []SectionLink // for the navbar
	// Data is for passing any page-specific data to the template
	Data any
}

type HomePageData struct {
	WritingEntries []*Entry
	MicroEntries   []*Entry
	Books          []BookReview
}

type BooksPageData struct {
	Books []BookReview
}

type ListPageData struct {
	Section string
	Entries []*Entry
}

type ShowPageData struct {
	Section string
	Slug    string
	Entry   *Entry
	// this is for books rn only but its basically additional data besides show page specific stuff
	Data any
}

type ArchivePageData struct {
	Section     string
	Tag         string
	ArchiveYear int
	Entries     []*Entry
}
