package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"app/store"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %q", resp["status"])
	}
}

func TestCreateBookmark(t *testing.T) {
	bookmarkStore = store.NewStore()

	bookmark := store.Bookmark{
		URL:         "https://example.com",
		Title:       "Example",
		Description: "Example description",
		Tags: []store.Tag{
			{Name: "web", Color: "#e0e0e0"},
		},
	}

	body, _ := json.Marshal(bookmark)
	req := httptest.NewRequest("POST", "/api/bookmarks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleCreateBookmark(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var created store.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if created.ID == "" {
		t.Error("Expected bookmark to have an ID")
	}
	if created.URL != bookmark.URL {
		t.Errorf("Expected URL %q, got %q", bookmark.URL, created.URL)
	}
	if created.Title != bookmark.Title {
		t.Errorf("Expected title %q, got %q", bookmark.Title, created.Title)
	}
}

func TestCreateBookmarkInvalidJSON(t *testing.T) {
	bookmarkStore = store.NewStore()

	req := httptest.NewRequest("POST", "/api/bookmarks", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleCreateBookmark(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListBookmarks(t *testing.T) {
	bookmarkStore = store.NewStore()

	b1 := &store.Bookmark{URL: "https://example1.com", Title: "Example 1"}
	b2 := &store.Bookmark{URL: "https://example2.com", Title: "Example 2"}

	bookmarkStore.CreateBookmark(b1)
	// Add a delay to ensure different IDs (generateID uses seconds precision)
	time.Sleep(1100 * time.Millisecond)
	bookmarkStore.CreateBookmark(b2)

	req := httptest.NewRequest("GET", "/api/bookmarks", nil)
	w := httptest.NewRecorder()

	handleListBookmarks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var bookmarks []*store.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&bookmarks); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(bookmarks))
	}
}

func TestListBookmarksEmpty(t *testing.T) {
	bookmarkStore = store.NewStore()

	req := httptest.NewRequest("GET", "/api/bookmarks", nil)
	w := httptest.NewRecorder()

	handleListBookmarks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var bookmarks []*store.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&bookmarks); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(bookmarks) != 0 {
		t.Errorf("Expected 0 bookmarks, got %d", len(bookmarks))
	}
}

func TestGetBookmark(t *testing.T) {
	bookmarkStore = store.NewStore()

	b := &store.Bookmark{URL: "https://example.com", Title: "Example"}
	bookmarkStore.CreateBookmark(b)

	req := httptest.NewRequest("GET", "/api/bookmarks/"+b.ID, nil)
	w := httptest.NewRecorder()

	handleGetBookmark(w, req, b.ID)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var retrieved store.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&retrieved); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if retrieved.ID != b.ID {
		t.Errorf("Expected ID %q, got %q", b.ID, retrieved.ID)
	}
}

func TestGetBookmarkNotFound(t *testing.T) {
	bookmarkStore = store.NewStore()

	req := httptest.NewRequest("GET", "/api/bookmarks/nonexistent", nil)
	w := httptest.NewRecorder()

	handleGetBookmark(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateBookmark(t *testing.T) {
	bookmarkStore = store.NewStore()

	b := &store.Bookmark{URL: "https://example.com", Title: "Original"}
	bookmarkStore.CreateBookmark(b)

	updated := store.Bookmark{
		URL:   "https://updated.com",
		Title: "Updated",
	}

	body, _ := json.Marshal(updated)
	req := httptest.NewRequest("PATCH", "/api/bookmarks/"+b.ID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleUpdateBookmark(w, req, b.ID)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result store.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.URL != updated.URL {
		t.Errorf("Expected URL %q, got %q", updated.URL, result.URL)
	}
	if result.Title != updated.Title {
		t.Errorf("Expected title %q, got %q", updated.Title, result.Title)
	}
}

func TestUpdateBookmarkNotFound(t *testing.T) {
	bookmarkStore = store.NewStore()

	updated := store.Bookmark{URL: "https://example.com", Title: "Updated"}
	body, _ := json.Marshal(updated)
	req := httptest.NewRequest("PATCH", "/api/bookmarks/nonexistent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleUpdateBookmark(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteBookmark(t *testing.T) {
	bookmarkStore = store.NewStore()

	b := &store.Bookmark{URL: "https://example.com", Title: "Example"}
	bookmarkStore.CreateBookmark(b)

	req := httptest.NewRequest("DELETE", "/api/bookmarks/"+b.ID, nil)
	w := httptest.NewRecorder()

	handleDeleteBookmark(w, req, b.ID)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify bookmark is deleted
	retrieved, err := bookmarkStore.GetBookmark(b.ID)
	if err == nil {
		t.Errorf("Expected bookmark to be deleted, but it still exists: %v", retrieved)
	}
}

func TestDeleteBookmarkNotFound(t *testing.T) {
	bookmarkStore = store.NewStore()

	req := httptest.NewRequest("DELETE", "/api/bookmarks/nonexistent", nil)
	w := httptest.NewRecorder()

	handleDeleteBookmark(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestSearchBookmarksByQuery(t *testing.T) {
	bookmarkStore = store.NewStore()

	b1 := &store.Bookmark{
		URL:   "https://example.com",
		Title: "Example",
		Description: "This is about golang",
	}
	b2 := &store.Bookmark{
		URL:   "https://other.com",
		Title: "Other",
		Description: "This is about python",
	}

	bookmarkStore.CreateBookmark(b1)
	time.Sleep(1100 * time.Millisecond)
	bookmarkStore.CreateBookmark(b2)

	filter := store.SearchFilter{Query: "golang"}
	body, _ := json.Marshal(filter)
	req := httptest.NewRequest("POST", "/api/bookmarks/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleSearchBookmarks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var results []*store.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestSearchBookmarksByTags(t *testing.T) {
	bookmarkStore = store.NewStore()

	b1 := &store.Bookmark{
		URL:   "https://example.com",
		Title: "Example",
		Tags:  []store.Tag{{Name: "go", Color: "#e0e0e0"}},
	}
	b2 := &store.Bookmark{
		URL:   "https://other.com",
		Title: "Other",
		Tags:  []store.Tag{{Name: "python", Color: "#e0e0e0"}},
	}

	bookmarkStore.CreateBookmark(b1)
	time.Sleep(1100 * time.Millisecond)
	bookmarkStore.CreateBookmark(b2)

	filter := store.SearchFilter{Tags: []string{"go"}}
	body, _ := json.Marshal(filter)
	req := httptest.NewRequest("POST", "/api/bookmarks/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleSearchBookmarks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var results []*store.Bookmark
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestSearchBookmarksInvalidJSON(t *testing.T) {
	bookmarkStore = store.NewStore()

	req := httptest.NewRequest("POST", "/api/bookmarks/search", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleSearchBookmarks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDashboardPage(t *testing.T) {
	bookmarkStore = store.NewStore()

	b := &store.Bookmark{
		URL:   "https://example.com",
		Title: "Example Site",
		Description: "An example site",
	}
	bookmarkStore.CreateBookmark(b)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected HTML content type, got %q", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()
	if !strings.Contains(body, "Bookmark Dashboard") {
		t.Error("Expected dashboard title in HTML")
	}
	if !strings.Contains(body, "Example Site") {
		t.Error("Expected bookmark title in HTML")
	}
}

func TestDashboardPageEmpty(t *testing.T) {
	bookmarkStore = store.NewStore()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "No bookmarks yet") {
		t.Error("Expected empty message in HTML")
	}
}

func TestCORSHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("Expected CORS origin *, got %q", origin)
	}

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if !strings.Contains(methods, "GET") {
		t.Errorf("Expected GET in CORS methods, got %q", methods)
	}
}

func TestOptionsRequestCORS(t *testing.T) {
	w := httptest.NewRecorder()

	// Simulate the OPTIONS handler in main
	setCORSHeaders(w)
	w.WriteHeader(http.StatusOK)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("Expected CORS origin *, got %q", origin)
	}
}

func TestIntegrationCreateAndRetrieve(t *testing.T) {
	bookmarkStore = store.NewStore()

	// Create
	bookmark := store.Bookmark{
		URL:   "https://example.com",
		Title: "Example",
	}
	body, _ := json.Marshal(bookmark)
	createReq := httptest.NewRequest("POST", "/api/bookmarks", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	handleCreateBookmark(createW, createReq)

	var created store.Bookmark
	json.NewDecoder(createW.Body).Decode(&created)

	// Retrieve
	getReq := httptest.NewRequest("GET", "/api/bookmarks/"+created.ID, nil)
	getW := httptest.NewRecorder()
	handleGetBookmark(getW, getReq, created.ID)

	var retrieved store.Bookmark
	json.NewDecoder(getW.Body).Decode(&retrieved)

	if retrieved.ID != created.ID {
		t.Errorf("Retrieved ID mismatch: %q vs %q", retrieved.ID, created.ID)
	}
	if retrieved.Title != created.Title {
		t.Errorf("Retrieved title mismatch: %q vs %q", retrieved.Title, created.Title)
	}
}

func TestIntegrationFullCycle(t *testing.T) {
	bookmarkStore = store.NewStore()

	// Create
	bookmark := store.Bookmark{
		URL:         "https://example.com",
		Title:       "Example",
		Description: "Test bookmark",
		Tags:        []store.Tag{{Name: "test", Color: "#e0e0e0"}},
	}
	body, _ := json.Marshal(bookmark)
	createReq := httptest.NewRequest("POST", "/api/bookmarks", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	handleCreateBookmark(createW, createReq)

	var created store.Bookmark
	json.NewDecoder(createW.Body).Decode(&created)
	createdID := created.ID

	// Update
	updated := store.Bookmark{
		URL:   "https://updated.com",
		Title: "Updated",
	}
	updateBody, _ := json.Marshal(updated)
	updateReq := httptest.NewRequest("PATCH", "/api/bookmarks/"+createdID, bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	handleUpdateBookmark(updateW, updateReq, createdID)

	if updateW.Code != http.StatusOK {
		t.Errorf("Update failed with status %d", updateW.Code)
	}

	// Delete
	deleteReq := httptest.NewRequest("DELETE", "/api/bookmarks/"+createdID, nil)
	deleteW := httptest.NewRecorder()
	handleDeleteBookmark(deleteW, deleteReq, createdID)

	if deleteW.Code != http.StatusNoContent {
		t.Errorf("Delete failed with status %d", deleteW.Code)
	}

	// Verify deleted
	getReq := httptest.NewRequest("GET", "/api/bookmarks/"+createdID, nil)
	getW := httptest.NewRecorder()
	handleGetBookmark(getW, getReq, createdID)

	if getW.Code != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", getW.Code)
	}
}
