package content

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"

	"new-portfolio/internal/models"

	"github.com/yuin/goldmark"
)

func sortEntries(entries []*models.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PublishedAt.After(entries[j].PublishedAt)
	})
}

func sortBooks(books []models.BookReview) {
	sort.Slice(books, func(i, j int) bool {
		return books[i].FinishedAt.After(books[j].FinishedAt)
	})
}

func newEntry(section, slug, title, summary string, publishedAt time.Time, tags []string) *models.Entry {
	body, err := mdFileToHtml(filepath.Join("./ui/mds/", section, slug+".md"))
	if err != nil {
		panic(err)
	}

	return &models.Entry{
		Slug:        slug,
		Title:       title,
		Summary:     summary,
		BodyHTML:    body,
		Section:     section,
		PublishedAt: publishedAt,
		Tags:        tags,
		// Comments: ,
	}
}

func newBookReview(title, slug, author, finishedAt string, rating int) models.BookReview {
	body, err := mdFileToHtml(filepath.Join("./ui/mds/", "books", slug+".md"))
	if err != nil {
		panic(err)
	}

	t, _ := time.Parse("2006-01-02", finishedAt)
	if rating < 0 {
		rating = 0
	}
	if rating > 5 {
		rating = 5
	}

	return models.BookReview{
		Entry: models.Entry{
			Slug:        slug,
			Title:       title,
			Summary:     "",
			BodyHTML:    body,
			Section:     "books",
			PublishedAt: t,
			Tags:        nil,
		},
		Author:     author,
		FinishedAt: t,
		Rating:     rating,
	}
}

func mdFileToHtml(pathToMd string) (template.HTML, error) {
	mdFile, err := os.ReadFile(pathToMd)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(mdFile, &buf); err != nil {
		return "", err
	}

	return template.HTML(buf.Bytes()), nil
}
