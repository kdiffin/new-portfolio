package content

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "new-portfolio/internal/db/sqlc"
	"new-portfolio/internal/models"
)

var errNotFound = errors.New("entry not found")

type Store struct {
	bySection map[string][]*models.Entry
	queries   *db.Queries // or whatever your queries type is

	books []models.BookReview
}

func NewStore(queries *db.Queries) *Store {
	books := []models.BookReview{
		newBookReview("A Philosophy of Software Design", "a-philosophy-of-software-design", "John Ousterhout", "2026-03-01", 5),
		// newBookReview("Apple in China", "a-philosophy-of-software-design", "Patrick McGee", "2026-02-11", 5),
		// newBookReview("Six Centuries of Type and Printing", "a-philosophy-of-software-design", "Glenn Fleishman", "2026-02-01", 4),
		// newBookReview("Pageboy", "a-philosophy-of-software-design", "Elliot Page", "2026-01-14", 4),
		// newBookReview("Every Day Is For The Thief", "a-philosophy-of-software-design", "Teju Cole", "2026-01-10", 5),
		// newBookReview("There Is No Antimemetics Division", "a-philosophy-of-software-design", "qntm", "2025-12-14", 2),
		// newBookReview("The Fort Bragg Cartel", "a-philosophy-of-software-design", "Seth Harp", "2025-12-02", 5),
		// newBookReview("Jesus and John Wayne", "a-philosophy-of-software-design", "Kristin Kobes Du Mez", "2025-11-27", 4),
		// newBookReview("Ways and Means", "a-philosophy-of-software-design", "Daniel Lefferts", "2025-11-15", 4),
		// newBookReview("Bad Monkey", "a-philosophy-of-software-design", "Carl Hiaasen", "2025-11-06", 3),
		// newBookReview("Timequake", "a-philosophy-of-software-design", "Kurt Vonnegut", "2025-10-15", 4),
		// newBookReview("Jewish Space Lasers", "a-philosophy-of-software-design", "Mike Rothschild", "2025-10-07", 3),
		// newBookReview("Uncommon Carriers", "a-philosophy-of-software-design", "John McPhee", "2025-09-05", 4),
		// newBookReview("Dawn", "a-philosophy-of-software-design", "Octavia E. Butler", "2025-08-23", 5),
		// newBookReview("Radio Free Albemuth", "a-philosophy-of-software-design", "Philip K Dick", "2025-07-27", 3),
		// newBookReview("Unruly", "a-philosophy-of-software-design", "David Mitchell", "2025-07-21", 3),
		// newBookReview("Things Become Other Things", "a-philosophy-of-software-design", "Craig Mod", "2025-06-29", 5),
		// newBookReview("Glass Century", "a-philosophy-of-software-design", "Ross Barkan", "2025-06-21", 5),
		// newBookReview("Abundance", "a-philosophy-of-software-design", "Ezra Klein and Derek Thompson", "2025-05-16", 3),
		// newBookReview("Careless People", "a-philosophy-of-software-design", "Sarah Wynn-Williams", "2025-04-06", 4),
		// newBookReview("Land is a Big Deal", "a-philosophy-of-software-design", "Lars A. Doucet", "2025-03-22", 5),
		// newBookReview("Cyberlibertarianism", "a-philosophy-of-software-design", "David Golumbia", "2025-02-23", 2),
		// newBookReview("Useful Not True", "a-philosophy-of-software-design", "Derek Sivers", "2025-02-15", 4),
		// newBookReview("The Hidden Wealth of Nations", "a-philosophy-of-software-design", "Gabriel Zucman", "2025-02-08", 4),
	}
	sortBooks(books)

	seed := map[string][]*models.Entry{
		"writing": {
			newEntry("writing", "2-types-of-work", "Software work is split into 2 parts", "cant be bothered to write a summary for seo", time.Date(2026, time.March, 28, 0, 0, 0, 0, time.UTC), []string{"intro", "site"}),
			newEntry("writing", "engineers-should-think", "To be an engineer, is to think, yet we're told not?", "Everything being measured via velocity and output has ruined our learning.", time.Date(2026, time.April, 26, 0, 0, 0, 0, time.UTC), []string{"micro"}),
		},

		"projects": {
			// newEntryWithQueries("projects", "portfolio-v1", "Portfolio v1", "Project entry placeholder.", now.AddDate(0, -1, 0), []string{"project", "web"}),
		},
		"micro": {
			newEntry("micro", "m-001", "hey user", "Short-form placeholder entry.", time.Date(2026, time.March, 28, 0, 0, 0, 0, time.UTC), []string{"micro"}),
			newEntry("micro", "m-002", "I don't like it when people use the phrase 'I don't get how someone can do x'", "just say you dont like it", time.Date(2026, time.March, 30, 0, 0, 0, 0, time.UTC), []string{"micro"}),

			newEntry("micro", "m-003", "chatgpt's way of speech when it knows you're a tech bro is so fucking annoying ", "dude when I ask how to manage my busy schedule there's no reason you should tie that to how that 'will make me a great AI engineer' sob emoji", time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC), []string{"micro"}),
			newEntry("micro", "m-004", "i like using the word orthogonal", "things in software shouldn't know about eachothers internals (hot take i know!!1)", time.Date(2026, time.April, 13, 0, 0, 0, 0, time.UTC), []string{"micro"}),
			newEntry("micro", "m-005", "How to know if someone's in AI psychosis", "control theory distributed system high signal buzzword.", time.Date(2026, time.April, 19, 0, 0, 0, 0, time.UTC), []string{"micro"}),
		},
		"books": {},
	}

	//  pass queries to newEntry

	for i := range books {
		entry := books[i].Entry
		seed["books"] = append(seed["books"], &entry)
	}

	for section := range seed {
		sortEntries(seed[section])
	}

	return &Store{bySection: seed, books: books, queries: queries}
}
func (s *Store) GetEntry(section, slug string) *models.Entry {
	entry, err := s.Find(section, slug)
	if err != nil {
		return nil
	}

	entry.Comments = s.commentsForEntry(section, slug)
	return entry
}

func (s *Store) commentsForEntry(section, slug string) []db.GetCommentsBySlugRow {
	if s.queries == nil {
		return nil
	}

	ctx := context.Background()
	pathSlug := fmt.Sprintf("/%s/%s", section, slug)

	rows, err := s.queries.GetCommentsBySlug(ctx, pathSlug)
	if err == nil && len(rows) > 0 {
		return rows
	}
	for i := range rows {
		fmt.Print(rows[i].ID)
		fmt.Print(rows[i].Name)
		fmt.Print(rows[i].Content)
		fmt.Print(rows[i].Slug)
		fmt.Print(rows[i].CreatedAt)
		fmt.Print(rows[i].ParentID)
	}

	fallbackRows, fallbackErr := s.queries.GetCommentsBySlug(ctx, slug)
	if fallbackErr == nil {
		return fallbackRows
	}

	return nil
}

func (s *Store) AddComment(section, slug, name, content string, parentID *int64) error {
	if s.queries == nil {
		return errors.New("comments datastore unavailable")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = "Anonymous"
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("comment content is required")
	}

	targetSlug := fmt.Sprintf("/%s/%s", section, slug)
	params := db.CreateCommentParams{
		Name:    name,
		Content: content,
		Slug:    targetSlug,
	}
	if parentID != nil {
		params.ParentID = sql.NullInt64{Int64: *parentID, Valid: true}
	}

	return s.queries.CreateComment(context.Background(), params)
}

func (s *Store) CommentExists(section, slug string, id int64) bool {
	if s.queries == nil {
		return false
	}

	comments := s.commentsForEntry(section, slug)
	for _, c := range comments {
		if c.ID == id {
			return true
		}
	}

	return false
}

func (s *Store) Sections() []models.SectionLink {
	return []models.SectionLink{
		{Label: "Writing", Href: "/writing"},
		{Label: "Books", Href: "/books"},
		{Label: "Projects", Href: "/projects"},
		{Label: "Micro", Href: "/micro"},
		{Label: "About", Href: "/about"},
	}
}
