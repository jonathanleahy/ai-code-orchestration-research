package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"task-board/model"
)

var store = model.NewStore()

func main() {
	http.HandleFunc("/api/tasks", handleTasks)
	http.HandleFunc("/api/tasks/", handleTask)
	http.HandleFunc("/", handleUI)

	port := "8890"
	log.Printf("task-board running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	switch r.Method {
	case http.MethodGet:
		status := r.URL.Query().Get("status")
		var tasks []*model.Task
		if status != "" {
			s := model.Status(strings.ToUpper(status))
			tasks = store.List(&s)
		} else {
			tasks = store.List(nil)
		}
		if tasks == nil {
			tasks = []*model.Task{}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})

	case http.MethodPost:
		var body struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}
		task, err := store.Create(body.Title, body.Description)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		http.Error(w, `{"error":"missing task ID"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := store.Get(id)
		if err != nil {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(task)

	case http.MethodPatch:
		var body struct {
			Title       *string       `json:"title"`
			Description *string       `json:"description"`
			Status      *model.Status `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}
		task, err := store.Update(id, body.Title, body.Description, body.Status)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(task)

	case http.MethodDelete:
		if err := store.Delete(id); err != nil {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"deleted": true})

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, uiHTML)
}

const uiHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Task Board</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f0f2f5; color: #1a1a2e; padding: 20px; }
h1 { text-align: center; margin-bottom: 24px; font-size: 24px; }
.board { display: flex; gap: 16px; max-width: 1200px; margin: 0 auto; }
.column { flex: 1; background: #fff; border-radius: 12px; padding: 16px; min-height: 400px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.column h2 { font-size: 14px; text-transform: uppercase; letter-spacing: 1px; color: #666; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 2px solid #eee; }
.column h2 span { background: #e8e8e8; border-radius: 12px; padding: 2px 8px; font-size: 12px; margin-left: 6px; }
.task { background: #fafafa; border: 1px solid #e8e8e8; border-radius: 8px; padding: 12px; margin-bottom: 8px; cursor: pointer; transition: box-shadow 0.2s; }
.task:hover { box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.task h3 { font-size: 14px; margin-bottom: 4px; }
.task p { font-size: 12px; color: #888; }
.task .actions { margin-top: 8px; display: flex; gap: 4px; }
.task button { font-size: 11px; padding: 2px 8px; border: 1px solid #ddd; border-radius: 4px; background: #fff; cursor: pointer; }
.task button:hover { background: #f0f0f0; }
.task button.delete { color: #e74c3c; border-color: #e74c3c; }
.add-form { margin-bottom: 16px; display: flex; gap: 8px; max-width: 1200px; margin: 0 auto 16px; }
.add-form input { flex: 1; padding: 8px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; }
.add-form button { padding: 8px 16px; background: #2ecc71; color: #fff; border: none; border-radius: 8px; cursor: pointer; font-size: 14px; }
.add-form button:hover { background: #27ae60; }
.empty { text-align: center; color: #ccc; padding: 40px; font-size: 14px; }
</style>
</head>
<body>
<h1>📋 Task Board</h1>

<div class="add-form">
  <input type="text" id="title" placeholder="New task title...">
  <input type="text" id="desc" placeholder="Description (optional)">
  <button onclick="addTask()">+ Add Task</button>
</div>

<div class="board">
  <div class="column" id="col-TODO"><h2>To Do <span id="count-TODO">0</span></h2><div id="tasks-TODO"></div></div>
  <div class="column" id="col-DOING"><h2>In Progress <span id="count-DOING">0</span></h2><div id="tasks-DOING"></div></div>
  <div class="column" id="col-DONE"><h2>Done <span id="count-DONE">0</span></h2><div id="tasks-DONE"></div></div>
</div>

<script>
const API = '/api/tasks';

async function load() {
  const res = await fetch(API);
  const data = await res.json();
  const tasks = data.tasks || [];
  ['TODO','DOING','DONE'].forEach(s => {
    const col = document.getElementById('tasks-' + s);
    const filtered = tasks.filter(t => t.status === s);
    document.getElementById('count-' + s).textContent = filtered.length;
    col.innerHTML = filtered.length ? filtered.map(t => taskHTML(t)).join('') : '<div class="empty">No tasks</div>';
  });
}

function taskHTML(t) {
  const next = {TODO:'DOING',DOING:'DONE',DONE:'TODO'};
  const prev = {DONE:'DOING',DOING:'TODO',TODO:'DONE'};
  return '<div class="task"><h3>' + esc(t.title) + '</h3>' +
    (t.description ? '<p>' + esc(t.description) + '</p>' : '') +
    '<p>Created: ' + new Date(t.createdAt).toLocaleString() + '</p>' +
    '<div class="actions">' +
    '<button onclick="move(\'' + t.id + '\',\'' + next[t.status] + '\')">→ ' + next[t.status] + '</button>' +
    '<button onclick="move(\'' + t.id + '\',\'' + prev[t.status] + '\')">← ' + prev[t.status] + '</button>' +
    '<button class="delete" onclick="del(\'' + t.id + '\')">✕</button>' +
    '</div></div>';
}

function esc(s) { const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

async function addTask() {
  const title = document.getElementById('title').value.trim();
  if (!title) return;
  const desc = document.getElementById('desc').value.trim();
  await fetch(API, { method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({title, description: desc}) });
  document.getElementById('title').value = '';
  document.getElementById('desc').value = '';
  load();
}

async function move(id, status) {
  await fetch(API + '/' + id, { method: 'PATCH', headers: {'Content-Type':'application/json'}, body: JSON.stringify({status}) });
  load();
}

async function del(id) {
  await fetch(API + '/' + id, { method: 'DELETE' });
  load();
}

load();
setInterval(load, 5000);
</script>
</body>
</html>`
