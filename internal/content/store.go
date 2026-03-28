package content

import (
	"errors"
	"slices"
	"sort"
	"strings"
	"time"

	"new-portfolio/internal/models"
)

var errNotFound = errors.New("entry not found")

type Store struct {
	bySection map[string][]*models.Entry
}

func NewStore() *Store {
	now := time.Now()
	seed := map[string][]*models.Entry{
		"writing": {
			newEntry("writing", "proto-note", "I haven't made anything with AT Proto yet", "Replace this with your first long-form article.", now.AddDate(0, 0, -6), []string{"intro", "site"}),
			newEntry("writing", "recently", "Recently", "A rolling updates post.", now.AddDate(0, 0, -21), []string{"notes"}),
			newEntry("writing", "new-tote", "New tote bag", "A short field note.", now.AddDate(0, 0, -26), []string{"notes"}),
			newEntry("writing", "color-dithering", "Color dithering", "A practical write-up.", now.AddDate(0, 0, -29), []string{"design"}),
			newEntry("writing", "sketches", "Sketches", "A collection of experiments.", now.AddDate(0, 1, -22), []string{"drawing"}),
			newEntry("writing", "theme-selector", "Theme selector", "How to implement theme preferences.", now.AddDate(0, 3, -14), []string{"site"}),
			newEntry("writing", "specifiers-history", "A brief history of specifiers and protocols", "Protocol and URL conventions.", now.AddDate(0, 3, -15), []string{"tech"}),
			newEntry("writing", "year-review", "2025 Year in Review", "Annual review and notes.", now.AddDate(0, 3, -16), []string{"review"}),
			newEntry("writing", "dec-recently", "Recently", "Another rolling post.", now.AddDate(0, 3, -22), []string{"notes"}),
		},
		"books-and-papers": {
			newEntry("books-and-papers", "book-log-1", "Reading Log 1", "Placeholder reading note.", now.AddDate(0, 0, -1), []string{"books"}),
		},
		"projects": {
			newEntry("projects", "portfolio-v1", "Portfolio v1", "Project entry placeholder.", now.AddDate(0, -1, 0), []string{"project", "web"}),
		},
		"micro": {
			newEntry("micro", "m-001", "Effect notes: PRs, progress, and joys", "Short-form placeholder entry.", now.AddDate(0, 0, -4), []string{"micro"}),
			newEntry("micro", "m-002", "Placemark & OSS changelog", "Another short-form placeholder.", now.AddDate(0, 0, -7), []string{"micro"}),
			newEntry("micro", "m-003", "ROOTS - Return Old Online Things to your own Site", "Another short-form placeholder.", now.AddDate(0, 0, -10), []string{"micro"}),
			newEntry("micro", "m-004", "Reactionary AI Centrism", "Another short-form placeholder.", now.AddDate(0, 1, -8), []string{"micro"}),
			newEntry("micro", "m-005", "Media diet", "Another short-form placeholder.", now.AddDate(0, 1, -15), []string{"micro"}),
		},
	}

	for section := range seed {
		sortEntries(seed[section])
	}

	return &Store{bySection: seed}
}

func (s *Store) Sections() []models.SectionLink {
	return []models.SectionLink{
		{Label: "Writing", Href: "/writing"},
		{Label: "Books and Papers", Href: "/books-and-papers"},
		{Label: "Projects", Href: "/projects"},
		{Label: "Micro", Href: "/micro"},
		{Label: "About", Href: "/about"},
	}
}

func (s *Store) List(section string) []*models.Entry {
	entries := slices.Clone(s.bySection[section])
	sortEntries(entries)
	return entries
}

func (s *Store) Find(section, slug string) (*models.Entry, error) {
	for _, entry := range s.bySection[section] {
		if entry.Slug == slug {
			return entry, nil
		}
	}
	return nil, errNotFound
}

func (s *Store) ListByYear(section string, year int) []*models.Entry {
	var out []*models.Entry
	for _, entry := range s.bySection[section] {
		if entry.PublishedAt.Year() == year {
			out = append(out, entry)
		}
	}
	sortEntries(out)
	return out
}

func (s *Store) ListByTag(section, tag string) []*models.Entry {
	var out []*models.Entry
	needle := strings.ToLower(tag)
	for _, entry := range s.bySection[section] {
		for _, t := range entry.Tags {
			if strings.EqualFold(t, needle) {
				out = append(out, entry)
				break
			}
		}
	}
	sortEntries(out)
	return out
}

func (s *Store) Latest(limit int) []*models.Entry {
	var all []*models.Entry
	for _, entries := range s.bySection {
		all = append(all, entries...)
	}
	sortEntries(all)
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return slices.Clone(all)
}

func sortEntries(entries []*models.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PublishedAt.After(entries[j].PublishedAt)
	})
}

func newEntry(section, slug, title, summary string, publishedAt time.Time, tags []string) *models.Entry {
	return &models.Entry{
		Section:     section,
		Slug:        slug,
		Title:       title,
		Summary:     summary,
		PublishedAt: publishedAt,
		Tags:        tags,
		BodyHTML:    "<p>Replace this placeholder body with your own original content.</p>",
	}
}
