package content

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"new-portfolio/internal/models"

	"github.com/yuin/goldmark"
)

var errNotFound = errors.New("entry not found")

type Store struct {
	bySection map[string][]*models.Entry
	books     []models.BookReview
}

func NewStore() *Store {
	books := []models.BookReview{
		newBookReview("A Philosophy of Software Design", "a-philosophy-of-software-design", "John Ousterhout", "2026-03-01", 5),
		newBookReview("Apple in China", "a-philosophy-of-software-design", "Patrick McGee", "2026-02-11", 5),
		newBookReview("Six Centuries of Type and Printing", "a-philosophy-of-software-design", "Glenn Fleishman", "2026-02-01", 4),
		newBookReview("Pageboy", "a-philosophy-of-software-design", "Elliot Page", "2026-01-14", 4),
		newBookReview("Every Day Is For The Thief", "a-philosophy-of-software-design", "Teju Cole", "2026-01-10", 5),
		newBookReview("There Is No Antimemetics Division", "a-philosophy-of-software-design", "qntm", "2025-12-14", 2),
		newBookReview("The Fort Bragg Cartel", "a-philosophy-of-software-design", "Seth Harp", "2025-12-02", 5),
		newBookReview("Jesus and John Wayne", "a-philosophy-of-software-design", "Kristin Kobes Du Mez", "2025-11-27", 4),
		newBookReview("Ways and Means", "a-philosophy-of-software-design", "Daniel Lefferts", "2025-11-15", 4),
		newBookReview("Bad Monkey", "a-philosophy-of-software-design", "Carl Hiaasen", "2025-11-06", 3),
		newBookReview("Timequake", "a-philosophy-of-software-design", "Kurt Vonnegut", "2025-10-15", 4),
		newBookReview("Jewish Space Lasers", "a-philosophy-of-software-design", "Mike Rothschild", "2025-10-07", 3),
		newBookReview("Uncommon Carriers", "a-philosophy-of-software-design", "John McPhee", "2025-09-05", 4),
		newBookReview("Dawn", "a-philosophy-of-software-design", "Octavia E. Butler", "2025-08-23", 5),
		newBookReview("Radio Free Albemuth", "a-philosophy-of-software-design", "Philip K Dick", "2025-07-27", 3),
		newBookReview("Unruly", "a-philosophy-of-software-design", "David Mitchell", "2025-07-21", 3),
		newBookReview("Things Become Other Things", "a-philosophy-of-software-design", "Craig Mod", "2025-06-29", 5),
		newBookReview("Glass Century", "a-philosophy-of-software-design", "Ross Barkan", "2025-06-21", 5),
		newBookReview("Abundance", "a-philosophy-of-software-design", "Ezra Klein and Derek Thompson", "2025-05-16", 3),
		newBookReview("Careless People", "a-philosophy-of-software-design", "Sarah Wynn-Williams", "2025-04-06", 4),
		newBookReview("Land is a Big Deal", "a-philosophy-of-software-design", "Lars A. Doucet", "2025-03-22", 5),
		newBookReview("Cyberlibertarianism", "a-philosophy-of-software-design", "David Golumbia", "2025-02-23", 2),
		newBookReview("Useful Not True", "a-philosophy-of-software-design", "Derek Sivers", "2025-02-15", 4),
		newBookReview("The Hidden Wealth of Nations", "a-philosophy-of-software-design", "Gabriel Zucman", "2025-02-08", 4),
	}
	sortBooks(books)

	seed := map[string][]*models.Entry{
		"writing": {
			newEntry("writing", "test", "I haven't made anything with AT Proto yet", "Replace this with your first long-form article.", time.Now().AddDate(0, 0, -6), []string{"intro", "site"}),
			// newEntry("writing", "recently", "Recently", "A rolling updates post.", now.AddDate(0, 0, -21), []string{"notes"}),
			// newEntry("writing", "new-tote", "New tote bag", "A short field note.", now.AddDate(0, 0, -26), []string{"notes"}),
			// newEntry("writing", "color-dithering", "Color dithering", "A practical write-up.", now.AddDate(0, 0, -29), []string{"design"}),
			// newEntry("writing", "sketches", "Sketches", "A collection of experiments.", now.AddDate(0, 1, -22), []string{"drawing"}),
			// newEntry("writing", "theme-selector", "Theme selector", "How to implement theme preferences.", now.AddDate(0, 3, -14), []string{"site"}),
			// newEntry("writing", "specifiers-history", "A brief history of specifiers and protocols", "Protocol and URL conventions.", now.AddDate(0, 3, -15), []string{"tech"}),
			// newEntry("writing", "year-review", "2025 Year in Review", "Annual review and notes.", now.AddDate(0, 3, -16), []string{"review"}),
			// newEntry("writing", "dec-recently", "Recently", "Another rolling post.", now.AddDate(0, 3, -22), []string{"notes"}),
		},

		"projects": {
			// newEntry("projects", "portfolio-v1", "Portfolio v1", "Project entry placeholder.", now.AddDate(0, -1, 0), []string{"project", "web"}),
		},
		"micro": {
			newEntry("micro", "m-001", "Effect notes: PRs, progress, and joys", "Short-form placeholder entry.", time.Now().AddDate(0, 0, -4), []string{"micro"}),
			// newEntry("micro", "m-002", "Placemark & OSS changelog", "Another short-form placeholder.", now.AddDate(0, 0, -7), []string{"micro"}),
			// newEntry("micro", "m-003", "ROOTS - Return Old Online Things to your own Site", "Another short-form placeholder.", now.AddDate(0, 0, -10), []string{"micro"}),
			// newEntry("micro", "m-004", "Reactionary AI Centrism", "Another short-form placeholder.", now.AddDate(0, 1, -8), []string{"micro"}),
			// newEntry("micro", "m-005", "Media diet", "Another short-form placeholder.", now.AddDate(0, 1, -15), []string{"micro"}),
		},
	}

	for section := range seed {
		sortEntries(seed[section])
	}

	return &Store{bySection: seed, books: books}
}

func (s *Store) Sections() []models.SectionLink {
	return []models.SectionLink{
		{Label: "Writing", Href: "/writing"},
		{Label: "Books", Href: "/books-and-papers"},
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

func (s *Store) Books() []models.BookReview {
	out := slices.Clone(s.books)
	sortBooks(out)
	return out
}

func (s *Store) LatestBooks(limit int) []models.BookReview {
	out := s.Books()
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

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
	var buf bytes.Buffer
	mdFile, err := os.ReadFile(filepath.Join("./ui/mds/", section, slug+".md"))
	if err != nil {
		panic(err)
	}

	if err := goldmark.Convert(mdFile, &buf); err != nil {
		panic(err)
	}

	print(buf.Bytes())

	return &models.Entry{
		Section:     section,
		Slug:        slug,
		Title:       title,
		Summary:     summary,
		PublishedAt: publishedAt,
		Tags:        tags,
		BodyHTML:    template.HTML(buf.String()),
	}
}

func newBookReview(title, slug, author, finishedAt string, rating int) models.BookReview {
	t, _ := time.Parse("2006-01-02", finishedAt)
	if rating < 0 {
		rating = 0
	}
	if rating > 5 {
		rating = 5
	}

	return models.BookReview{
		Title:      title,
		Author:     author,
		Slug:       slug,
		FinishedAt: t,
		Rating:     rating,
	}
}
