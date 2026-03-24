package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-board/model"
)

func setupServer() *http.ServeMux {
	store = model.NewStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", handleTasks)
	mux.HandleFunc("/api/tasks/", handleTask)
	mux.HandleFunc("/", handleUI)
	return mux
}

func TestUIReturnsHTML(t *testing.T) {
	mux := setupServer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("expected text/html, got %s", ct)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Task Board")) {
		t.Fatal("expected 'Task Board' in HTML")
	}
}

func TestCreateTaskViaAPI(t *testing.T) {
	mux := setupServer()
	body := `{"title":"Test Task","description":"A test"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var task model.Task
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Fatal(err)
	}
	if task.Title != "Test Task" {
		t.Errorf("expected title 'Test Task', got %q", task.Title)
	}
	if task.Status != model.StatusTodo {
		t.Errorf("expected TODO, got %q", task.Status)
	}
}

func TestCreateTaskEmptyTitle(t *testing.T) {
	mux := setupServer()
	body := `{"title":"","description":""}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListTasks(t *testing.T) {
	mux := setupServer()

	// Create two tasks
	for _, title := range []string{"Task 1", "Task 2"} {
		body := `{"title":"` + title + `"}`
		req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}

	// List all
	req := httptest.NewRequest("GET", "/api/tasks", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result struct {
		Tasks []model.Task `json:"tasks"`
	}
	json.NewDecoder(w.Body).Decode(&result)
	if len(result.Tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(result.Tasks))
	}
}

func TestListTasksWithFilter(t *testing.T) {
	mux := setupServer()

	// Create a task
	body := `{"title":"Filtered Task"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// List with filter TODO
	req = httptest.NewRequest("GET", "/api/tasks?status=todo", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var result struct {
		Tasks []model.Task `json:"tasks"`
	}
	json.NewDecoder(w.Body).Decode(&result)
	if len(result.Tasks) != 1 {
		t.Fatalf("expected 1 TODO task, got %d", len(result.Tasks))
	}

	// List with filter DOING (none)
	req = httptest.NewRequest("GET", "/api/tasks?status=doing", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	result.Tasks = nil
	json.NewDecoder(w.Body).Decode(&result)
	if len(result.Tasks) != 0 {
		t.Fatalf("expected 0 DOING tasks, got %d", len(result.Tasks))
	}
}

func TestGetTask(t *testing.T) {
	mux := setupServer()

	// Create
	body := `{"title":"Get Me"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)

	// Get
	req = httptest.NewRequest("GET", "/api/tasks/"+created.ID, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got model.Task
	json.NewDecoder(w.Body).Decode(&got)
	if got.Title != "Get Me" {
		t.Errorf("expected 'Get Me', got %q", got.Title)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	mux := setupServer()
	req := httptest.NewRequest("GET", "/api/tasks/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	mux := setupServer()

	// Create
	body := `{"title":"Move Me"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)

	// Update status to DOING
	body = `{"status":"DOING"}`
	req = httptest.NewRequest("PATCH", "/api/tasks/"+created.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated model.Task
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Status != model.StatusDoing {
		t.Errorf("expected DOING, got %q", updated.Status)
	}
}

func TestUpdateTaskTitle(t *testing.T) {
	mux := setupServer()

	body := `{"title":"Original"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)

	body = `{"title":"Updated"}`
	req = httptest.NewRequest("PATCH", "/api/tasks/"+created.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var updated model.Task
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Title != "Updated" {
		t.Errorf("expected 'Updated', got %q", updated.Title)
	}
}

func TestDeleteTask(t *testing.T) {
	mux := setupServer()

	body := `{"title":"Delete Me"}`
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)

	// Delete
	req = httptest.NewRequest("DELETE", "/api/tasks/"+created.ID, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Verify gone
	req = httptest.NewRequest("GET", "/api/tasks/"+created.ID, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	mux := setupServer()
	req := httptest.NewRequest("DELETE", "/api/tasks/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestEmptyListReturnsEmptyArray(t *testing.T) {
	mux := setupServer()
	req := httptest.NewRequest("GET", "/api/tasks", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var result struct {
		Tasks []model.Task `json:"tasks"`
	}
	json.NewDecoder(w.Body).Decode(&result)
	if result.Tasks == nil {
		t.Fatal("expected empty array, got nil")
	}
	if len(result.Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(result.Tasks))
	}
}
