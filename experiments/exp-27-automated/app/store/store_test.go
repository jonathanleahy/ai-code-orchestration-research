package store

import (
	"os"
	"reflect"
	"testing"
)

func TestCreateBookmark(t *testing.T) {
	store := NewStore()

	b := &Bookmark{
		URL:         "https://example.com",
		Title:       "Example",
		Description: "An example site",
		Tags: []Tag{
			{Name: "test", Color: "blue"},
		},
	}

	err := store.CreateBookmark(b)
	if err != nil {
		t.Fatalf("CreateBookmark failed: %v", err)
	}

	if b.ID == "" {
		t.Error("Bookmark ID should be generated")
	}

	if b.CreatedAt.IsZero() {
		t.Error("Bookmark CreatedAt should be set")
	}

	if b.UpdatedAt.IsZero() {
		t.Error("Bookmark UpdatedAt should be set")
	}

	// Verify the bookmark was actually stored
	retrieved, err := store.GetBookmark(b.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve created bookmark: %v", err)
	}

	if retrieved.ID != b.ID {
		t.Errorf("Expected ID %s, got %s", b.ID, retrieved.ID)
	}
}

func TestGetBookmark(t *testing.T) {
	store := NewStore()

	// Create a bookmark first
	b := &Bookmark{
		URL:   "https://example.com",
		Title: "Example",
	}
	store.CreateBookmark(b)

	// Test retrieving existing bookmark
	retrieved, err := store.GetBookmark(b.ID)
	if err != nil {
		t.Fatalf("GetBookmark failed: %v", err)
	}

	if retrieved.ID != b.ID {
		t.Errorf("Expected ID %s, got %s", b.ID, retrieved.ID)
	}

	// Test retrieving non-existing bookmark
	_, err = store.GetBookmark("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent bookmark")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Expected IsNotExist error, got %v", err)
	}
}

func TestUpdateBookmark(t *testing.T) {
	store := NewStore()

	// Create a bookmark first
	b := &Bookmark{
		URL:         "https://example.com",
		Title:       "Example",
		Description: "An example site",
		Tags: []Tag{
			{Name: "test", Color: "blue"},
		},
	}
	store.CreateBookmark(b)

	// Update the bookmark
	updated := &Bookmark{
		URL:         "https://updated.com",
		Title:       "Updated Example",
		Description: "An updated example site",
		Tags: []Tag{
			{Name: "updated", Color: "red"},
		},
	}

	err := store.UpdateBookmark(b.ID, updated)
	if err != nil {
		t.Fatalf("UpdateBookmark failed: %v", err)
	}

	// Verify the update
	retrieved, err := store.GetBookmark(b.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated bookmark: %v", err)
	}

	if retrieved.URL != updated.URL {
		t.Errorf("Expected URL %s, got %s", updated.URL, retrieved.URL)
	}

	if retrieved.Title != updated.Title {
		t.Errorf("Expected title %s, got %s", updated.Title, retrieved.Title)
	}

	if retrieved.Description != updated.Description {
		t.Errorf("Expected description %s, got %s", updated.Description, retrieved.Description)
	}

	if !reflect.DeepEqual(retrieved.Tags, updated.Tags) {
		t.Errorf("Expected tags %v, got %v", updated.Tags, retrieved.Tags)
	}

	if retrieved.UpdatedAt.IsZero() {
		t.Error("Bookmark UpdatedAt should be updated")
	}

	// Test updating non-existing bookmark
	err = store.UpdateBookmark("nonexistent", updated)
	if err == nil {
		t.Error("Expected error for non-existent bookmark")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Expected IsNotExist error, got %v", err)
	}
}

func TestDeleteBookmark(t *testing.T) {
	store := NewStore()

	// Create a bookmark first
	b := &Bookmark{
		URL:   "https://example.com",
		Title: "Example",
	}
	store.CreateBookmark(b)

	// Verify it exists
	_, err := store.GetBookmark(b.ID)
	if err != nil {
		t.Fatalf("Bookmark should exist: %v", err)
	}

	// Delete the bookmark
	err = store.DeleteBookmark(b.ID)
	if err != nil {
		t.Fatalf("DeleteBookmark failed: %v", err)
	}

	// Verify it's gone
	_, err = store.GetBookmark(b.ID)
	if err == nil {
		t.Error("Bookmark should be deleted")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Expected IsNotExist error, got %v", err)
	}

	// Test deleting non-existing bookmark
	err = store.DeleteBookmark("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent bookmark")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Expected IsNotExist error, got %v", err)
	}
}

func TestListBookmarks(t *testing.T) {
	store := NewStore()

	// Create some bookmarks
	b1 := &Bookmark{
		URL:   "https://example1.com",
		Title: "Example 1",
	}
	b2 := &Bookmark{
		URL:   "https://example2.com",
		Title: "Example 2",
	}
	store.CreateBookmark(b1)
	store.CreateBookmark(b2)

	// List all bookmarks
	bookmarks, err := store.ListBookmarks()
	if err != nil {
		t.Fatalf("ListBookmarks failed: %v", err)
	}

	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(bookmarks))
	}

	// Verify we got the right bookmarks
	foundB1 := false
	foundB2 := false
	for _, b := range bookmarks {
		if b.ID == b1.ID {
			foundB1 = true
		}
		if b.ID == b2.ID {
			foundB2 = true
		}
	}

	if !foundB1 {
		t.Error("Bookmark 1 not found in list")
	}
	if !foundB2 {
		t.Error("Bookmark 2 not found in list")
	}
}

func TestSearchBookmarks(t *testing.T) {
	store := NewStore()

	// Create some bookmarks with different tags
	b1 := &Bookmark{
		URL:         "https://example1.com",
		Title:       "Example 1",
		Description: "First example",
		Tags: []Tag{
			{Name: "tag1", Color: "red"},
			{Name: "tag2", Color: "blue"},
		},
	}
	b2 := &Bookmark{
		URL:         "https://example2.com",
		Title:       "Example 2",
		Description: "Second example",
		Tags: []Tag{
			{Name: "tag2", Color: "blue"},
			{Name: "tag3", Color: "green"},
		},
	}
	store.CreateBookmark(b1)
	store.CreateBookmark(b2)

	// Search by tag
	filter := SearchFilter{
		Query: "",
		Tags:  []string{"tag2"},
	}
	bookmarks, err := store.SearchBookmarks(filter)
	if err != nil {
		t.Fatalf("SearchBookmarks failed: %v", err)
	}

	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks with tag2, got %d", len(bookmarks))
	}

	// Search by query
	filter = SearchFilter{
		Query: "Example 1",
		Tags:  []string{},
	}
	bookmarks, err = store.SearchBookmarks(filter)
	if err != nil {
		t.Fatalf("SearchBookmarks failed: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark with query 'Example 1', got %d", len(bookmarks))
	}

	// Search by both query and tag
	filter = SearchFilter{
		Query: "Example",
		Tags:  []string{"tag1"},
	}
	bookmarks, err = store.SearchBookmarks(filter)
	if err != nil {
		t.Fatalf("SearchBookmarks failed: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark with query 'Example' and tag 'tag1', got %d", len(bookmarks))
	}
}
