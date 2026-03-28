package models

import (
	"html/template"
	"time"
)

type Entry struct {
	Slug        string
	Title       string
	Summary     string
	BodyHTML    template.HTML
	Section     string
	PublishedAt time.Time
	Tags        []string
}

type PageTemplateData struct {
	Title       string // for seo
	Description string // for seo
	Path        string // for the navbar highlighting
	Theme       string
	Year        int           // for the footer
	Sections    []SectionLink // for the navbar
	Data        any
}

type HomePageData struct {
	WritingEntries []*Entry
	MicroEntries   []*Entry
	Books          []BookReview
}

type BookReview struct {
	Entry
	Author     string
	FinishedAt time.Time
	Rating     int
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
	Book    *BookReview
}

type ArchivePageData struct {
	Section     string
	Tag         string
	ArchiveYear int
	Entries     []*Entry
}

type SectionLink struct {
	Label string
	Href  string
}
