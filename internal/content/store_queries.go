package content

import (
	"slices"
	"strings"

	"new-portfolio/internal/models"
)

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

func (s *Store) FindForShow(section, slug string) (models.ShowPageData, error) {
	if section == "books" {
		for i := range s.books {
			if s.books[i].Slug == slug {
				s.books[i].Entry.Comments = s.commentsForEntry(section, slug)
				return models.ShowPageData{
					Section: section,
					Slug:    slug,
					Entry:   &s.books[i].Entry,
					Data:    &s.books[i],
				}, nil
			}
		}
		return models.ShowPageData{}, errNotFound
	}

	entry, err := s.Find(section, slug)
	if err != nil {
		return models.ShowPageData{}, err
	}

	entry.Comments = s.commentsForEntry(section, slug)

	return models.ShowPageData{
		Section: section,
		Slug:    slug,
		Entry:   entry,
	}, nil
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
