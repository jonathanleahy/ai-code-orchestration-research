package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"app/store"
)

var bookmarkStore store.Store

func init() {
	bookmarkStore = store.NewStore()
}

func main() {
	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/health", handleHealth)

	// Bookmark CRUD endpoints
	http.HandleFunc("/api/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPost {
			handleCreateBookmark(w, r)
		} else if r.Method == http.MethodGet {
			handleListBookmarks(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/bookmarks/", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/bookmarks/")
		if r.Method == http.MethodGet {
			handleGetBookmark(w, r, id)
		} else if r.Method == http.MethodPatch {
			handleUpdateBookmark(w, r, id)
		} else if r.Method == http.MethodDelete {
			handleDeleteBookmark(w, r, id)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/bookmarks/search", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPost {
			handleSearchBookmarks(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	respondJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func handleCreateBookmark(w http.ResponseWriter, r *http.Request) {
	var b store.Bookmark
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bookmarkStore.CreateBookmark(&b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, b, http.StatusCreated)
}

func handleGetBookmark(w http.ResponseWriter, r *http.Request, id string) {
	b, err := bookmarkStore.GetBookmark(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	respondJSON(w, b, http.StatusOK)
}

func handleListBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := bookmarkStore.ListBookmarks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, bookmarks, http.StatusOK)
}

func handleUpdateBookmark(w http.ResponseWriter, r *http.Request, id string) {
	var b store.Bookmark
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bookmarkStore.UpdateBookmark(id, &b); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	b.ID = id
	respondJSON(w, b, http.StatusOK)
}

func handleDeleteBookmark(w http.ResponseWriter, r *http.Request, id string) {
	if err := bookmarkStore.DeleteBookmark(id); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleSearchBookmarks(w http.ResponseWriter, r *http.Request) {
	var filter store.SearchFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bookmarks, err := bookmarkStore.SearchBookmarks(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, bookmarks, http.StatusOK)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	bookmarks, _ := bookmarkStore.ListBookmarks()

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	html = html + "<meta charset=\"UTF-8\">\n"
	html = html + "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"
	html = html + "<title>Bookmark Dashboard</title>\n"
	html = html + "<style>\n"
	html = html + "body { font-family: system-ui, -apple-system, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }\n"
	html = html + ".container { max-width: 1200px; margin: 0 auto; }\n"
	html = html + ".header { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }\n"
	html = html + "h1 { margin: 0; color: #333; }\n"
	html = html + ".form-section { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }\n"
	html = html + ".form-group { margin-bottom: 15px; }\n"
	html = html + "label { display: block; font-weight: 600; margin-bottom: 5px; color: #333; }\n"
	html = html + "input, textarea { width: 100%; padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; font-family: inherit; box-sizing: border-box; }\n"
	html = html + "textarea { resize: vertical; min-height: 80px; }\n"
	html = html + "button { background: #0066cc; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-weight: 600; }\n"
	html = html + "button:hover { background: #0052a3; }\n"
	html = html + ".bookmarks-section { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }\n"
	html = html + ".bookmark-item { border: 1px solid #ddd; padding: 15px; border-radius: 4px; margin-bottom: 15px; background: #fafafa; }\n"
	html = html + ".bookmark-title { font-weight: 600; color: #0066cc; margin-bottom: 5px; }\n"
	html = html + ".bookmark-url { font-size: 14px; color: #666; margin-bottom: 8px; word-break: break-all; }\n"
	html = html + ".bookmark-desc { font-size: 14px; color: #444; margin-bottom: 10px; }\n"
	html = html + ".bookmark-tags { margin-bottom: 10px; }\n"
	html = html + ".tag { display: inline-block; padding: 4px 8px; border-radius: 3px; font-size: 12px; margin-right: 5px; margin-bottom: 5px; }\n"
	html = html + ".bookmark-actions { display: flex; gap: 10px; }\n"
	html = html + ".btn-small { padding: 6px 12px; font-size: 12px; }\n"
	html = html + ".btn-delete { background: #dc3545; }\n"
	html = html + ".btn-delete:hover { background: #c82333; }\n"
	html = html + ".empty { color: #999; font-style: italic; text-align: center; padding: 40px 20px; }\n"
	html = html + "</style>\n"
	html = html + "</head>\n<body>\n"
	html = html + "<div class=\"container\">\n"
	html = html + "<div class=\"header\"><h1>📚 Bookmark Dashboard</h1></div>\n"
	html = html + "<div class=\"form-section\">\n"
	html = html + "<h2>Add New Bookmark</h2>\n"
	html = html + "<form id=\"bookmarkForm\">\n"
	html = html + "<div class=\"form-group\">\n"
	html = html + "<label>URL</label>\n"
	html = html + "<input type=\"url\" id=\"url\" required placeholder=\"https://example.com\">\n"
	html = html + "</div>\n"
	html = html + "<div class=\"form-group\">\n"
	html = html + "<label>Title</label>\n"
	html = html + "<input type=\"text\" id=\"title\" required placeholder=\"Page title\">\n"
	html = html + "</div>\n"
	html = html + "<div class=\"form-group\">\n"
	html = html + "<label>Description</label>\n"
	html = html + "<textarea id=\"description\" placeholder=\"Optional description\"></textarea>\n"
	html = html + "</div>\n"
	html = html + "<div class=\"form-group\">\n"
	html = html + "<label>Tags (comma-separated)</label>\n"
	html = html + "<input type=\"text\" id=\"tags\" placeholder=\"web, tutorial, javascript\">\n"
	html = html + "</div>\n"
	html = html + "<button type=\"submit\">Add Bookmark</button>\n"
	html = html + "</form>\n"
	html = html + "</div>\n"
	html = html + "<div class=\"bookmarks-section\">\n"
	html = html + "<h2>Bookmarks (" + formatInt(len(bookmarks)) + ")</h2>\n"

	if len(bookmarks) == 0 {
		html = html + "<div class=\"empty\">No bookmarks yet. Add your first bookmark above!</div>\n"
	} else {
		for _, b := range bookmarks {
			html = html + "<div class=\"bookmark-item\">\n"
			html = html + "<div class=\"bookmark-title\">" + escapeHTML(b.Title) + "</div>\n"
			html = html + "<div class=\"bookmark-url\"><a href=\"" + escapeHTML(b.URL) + "\" target=\"_blank\">" + escapeHTML(b.URL) + "</a></div>\n"
			if b.Description != "" {
				html = html + "<div class=\"bookmark-desc\">" + escapeHTML(b.Description) + "</div>\n"
			}
			if len(b.Tags) > 0 {
				html = html + "<div class=\"bookmark-tags\">\n"
				for _, tag := range b.Tags {
					color := tag.Color
					if color == "" {
						color = "#e0e0e0"
					}
					html = html + "<span class=\"tag\" style=\"background: " + escapeHTML(color) + ";\">" + escapeHTML(tag.Name) + "</span>\n"
				}
				html = html + "</div>\n"
			}
			html = html + "<div class=\"bookmark-actions\">\n"
			html = html + "<button class=\"btn-small\" onclick=\"editBookmark('" + escapeHTML(b.ID) + "')\">Edit</button>\n"
			html = html + "<button class=\"btn-small btn-delete\" onclick=\"deleteBookmark('" + escapeHTML(b.ID) + "')\">Delete</button>\n"
			html = html + "</div>\n"
			html = html + "</div>\n"
		}
	}

	html = html + "</div>\n"
	html = html + "</div>\n"
	html = html + "<script>\n"
	html = html + "document.getElementById('bookmarkForm').addEventListener('submit', function(e) {\n"
	html = html + "  e.preventDefault();\n"
	html = html + "  const url = document.getElementById('url').value;\n"
	html = html + "  const title = document.getElementById('title').value;\n"
	html = html + "  const description = document.getElementById('description').value;\n"
	html = html + "  const tagsInput = document.getElementById('tags').value;\n"
	html = html + "  const tagsArray = tagsInput.split(',').map(t => t.trim()).filter(t => t !== '').map(t => ({name: t, color: '#e0e0e0'}));\n"
	html = html + "  const bookmark = {url: url, title: title, description: description, tags: tagsArray};\n"
	html = html + "  fetch('/api/bookmarks', {\n"
	html = html + "    method: 'POST',\n"
	html = html + "    headers: {'Content-Type': 'application/json'},\n"
	html = html + "    body: JSON.stringify(bookmark)\n"
	html = html + "  }).then(r => r.json()).then(data => {location.reload();}).catch(err => alert('Error: ' + err));\n"
	html = html + "});\n"
	html = html + "function deleteBookmark(id) {\n"
	html = html + "  if (!confirm('Delete this bookmark?')) return;\n"
	html = html + "  fetch('/api/bookmarks/' + id, {method: 'DELETE'}).then(() => {location.reload();}).catch(err => alert('Error: ' + err));\n"
	html = html + "}\n"
	html = html + "function editBookmark(id) {\n"
	html = html + "  alert('Edit functionality would be implemented here for bookmark: ' + id);\n"
	html = html + "}\n"
	html = html + "</script>\n"
	html = html + "</body>\n</html>\n"

	io.WriteString(w, html)
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

func formatInt(n int) string {
	switch n {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	default:
		return "many"
	}
}
