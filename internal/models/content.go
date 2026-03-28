package models

import (
	"html/template"
	"time"
)

type Entry struct {
	Section     string
	Slug        string
	Title       string
	Summary     string
	PublishedAt time.Time
	Tags        []string
	BodyHTML    template.HTML
}

type TemplateData struct {
	Title       string
	Description string
	Path        string
	Theme       string
	Year        int
	Sections    []SectionLink
	Data        any
}

type HomePageData struct {
	WritingEntries []*Entry
	MicroEntries   []*Entry
}

type ListPageData struct {
	Section string
	Entries []*Entry
}

type ShowPageData struct {
	Section string
	Slug    string
	Entry   *Entry
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
