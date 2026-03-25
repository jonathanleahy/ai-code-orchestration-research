package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

// Bookmark represents a single bookmark entry
type Bookmark struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []Tag     `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag represents a tag with color coding
type Tag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// SearchFilter holds parameters for filtering and searching
type SearchFilter struct {
	Query string   `json:"query"`
	Tags  []string `json:"tags"`
}

// ExportData holds the structure for exporting bookmarks
type ExportData struct {
	Bookmarks []Bookmark `json:"bookmarks"`
}

// Store interface defines the required methods for bookmark management
type Store interface {
	CreateBookmark(b *Bookmark) error
	GetBookmark(id string) (*Bookmark, error)
	UpdateBookmark(id string, b *Bookmark) error
	DeleteBookmark(id string) error
	ListBookmarks() ([]*Bookmark, error)
	SearchBookmarks(filter SearchFilter) ([]*Bookmark, error)
	FilterByTags(tags []string) ([]*Bookmark, error)
	Export() (*ExportData, error)
	Import(data *ExportData) error
	SaveToFile(filename string) error
	LoadFromFile(filename string) error
}

// store implements the Store interface
type store struct {
	bookmarks map[string]*Bookmark
	mu        sync.RWMutex
}

// NewStore creates and returns a new Store instance
func NewStore() Store {
	return &store{
		bookmarks: make(map[string]*Bookmark),
	}
}

// CreateBookmark adds a new bookmark to the store
func (s *store) CreateBookmark(b *Bookmark) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b.ID = generateID()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	s.bookmarks[b.ID] = b
	return nil
}

// GetBookmark retrieves a bookmark by its ID
func (s *store) GetBookmark(id string) (*Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, exists := s.bookmarks[id]
	if !exists {
		return nil, os.ErrNotExist
	}
	return b, nil
}

// UpdateBookmark updates an existing bookmark
func (s *store) UpdateBookmark(id string, b *Bookmark) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.bookmarks[id]
	if !exists {
		return os.ErrNotExist
	}

	existing.URL = b.URL
	existing.Title = b.Title
	existing.Description = b.Description
	existing.Tags = b.Tags
	existing.UpdatedAt = time.Now()
	return nil
}

// DeleteBookmark removes a bookmark by its ID
func (s *store) DeleteBookmark(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.bookmarks[id]; !exists {
		return os.ErrNotExist
	}
	delete(s.bookmarks, id)
	return nil
}

// ListBookmarks returns all bookmarks
func (s *store) ListBookmarks() ([]*Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bookmarks := make([]*Bookmark, 0, len(s.bookmarks))
	for _, b := range s.bookmarks {
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

// SearchBookmarks filters bookmarks based on query and tags
func (s *store) SearchBookmarks(filter SearchFilter) ([]*Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*Bookmark

	for _, b := range s.bookmarks {
		if filter.Query != "" {
			if !contains(b.Title, filter.Query) && !contains(b.Description, filter.Query) && !contains(b.URL, filter.Query) {
				continue
			}
		}

		if len(filter.Tags) > 0 {
			if !hasAnyTag(b.Tags, filter.Tags) {
				continue
			}
		}

		results = append(results, b)
	}

	return results, nil
}

// FilterByTags returns bookmarks that match any of the given tags
func (s *store) FilterByTags(tags []string) ([]*Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*Bookmark

	for _, b := range s.bookmarks {
		if hasAnyTag(b.Tags, tags) {
			results = append(results, b)
		}
	}

	return results, nil
}

// Export returns all bookmarks in export format
func (s *store) Export() (*ExportData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bookmarks := make([]Bookmark, 0, len(s.bookmarks))
	for _, b := range s.bookmarks {
		bookmarks = append(bookmarks, *b)
	}

	return &ExportData{Bookmarks: bookmarks}, nil
}

// Import loads bookmarks from export data
func (s *store) Import(data *ExportData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bookmarks = make(map[string]*Bookmark)
	for _, b := range data.Bookmarks {
		bookmark := b
		bookmark.CreatedAt = time.Now()
		bookmark.UpdatedAt = time.Now()
		s.bookmarks[bookmark.ID] = &bookmark
	}
	return nil
}

// SaveToFile saves all bookmarks to a JSON file
func (s *store) SaveToFile(filename string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.bookmarks, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// LoadFromFile loads bookmarks from a JSON file
func (s *store) LoadFromFile(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var bookmarks map[string]*Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return err
	}

	s.bookmarks = bookmarks
	return nil
}

// Helper function to check if a string contains another string (case insensitive)
func contains(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// Helper function to check if a bookmark has any of the specified tags
func hasAnyTag(bookmarkTags []Tag, filterTags []string) bool {
	for _, tag := range bookmarkTags {
		for _, filterTag := range filterTags {
			if strings.ToLower(tag.Name) == strings.ToLower(filterTag) {
				return true
			}
		}
	}
	return false
}

// Helper function to generate a unique ID
func generateID() string {
	return time.Now().Format("20060102150405")
}
