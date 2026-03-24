# task-board — Architecture Specification

## Overview
A Go backend with an in-memory task store. Contract-first: schema.graphql defines all types and operations.

## Contract (schema.graphql)
```graphql
type Task {
  id: ID!
  title: String!
  description: String
  status: Status!
  createdAt: String!
  updatedAt: String!
}

enum Status { TODO, DOING, DONE }

type Query {
  tasks(status: Status): [Task!]!
  task(id: ID!): Task
}

type Mutation {
  createTask(title: String!, description: String): Task!
  updateTask(id: ID!, title: String, description: String, status: Status): Task!
  deleteTask(id: ID!): Boolean!
}
```

## Files (Go only — no frontend in this experiment)

```
task-board/
  go.mod                    # module task-board, go 1.22
  schema.graphql            # The contract (above)
  model/task.go             # Task struct + in-memory Store (~100 lines)
  model/task_test.go        # 10 tests for Store CRUD (~90 lines)
```

## Module: model/task.go

**Package:** `model`

**Types:**
```go
type Status string
const (
    StatusTodo  Status = "TODO"
    StatusDoing Status = "DOING"
    StatusDone  Status = "DONE"
)

type Task struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      Status `json:"status"`
    CreatedAt   string `json:"createdAt"`
    UpdatedAt   string `json:"updatedAt"`
}
```

**Store:**
```go
type Store struct { /* mutex + map[string]*Task + nextID counter */ }

func NewStore() *Store
func (s *Store) Create(title, description string) (*Task, error)  // Returns error if title empty
func (s *Store) Get(id string) (*Task, error)                      // Returns error if not found
func (s *Store) List(status *Status) []*Task                       // nil status = all tasks
func (s *Store) Update(id string, title, description *string, status *Status) (*Task, error)
func (s *Store) Delete(id string) error                             // Returns error if not found
```

**Implementation notes:**
- Use `sync.RWMutex` for thread safety
- Auto-increment string IDs ("1", "2", "3"...)
- CreatedAt/UpdatedAt use `time.Now().UTC().Format(time.RFC3339)`
- New tasks default to `StatusTodo`
- `fmt.Errorf` for errors (no custom error types)

**Validation gate:** `go vet ./model/...` exits 0

## Module: model/task_test.go

**Package:** `model` (same package — white-box testing)

**10 tests:**
1. `TestCreate` — creates task, checks ID non-empty, title correct, status TODO
2. `TestCreateEmptyTitle` — empty title returns error
3. `TestGet` — creates then gets, title matches
4. `TestGetNotFound` — nonexistent ID returns error
5. `TestList` — creates 2 tasks, List(nil) returns 2
6. `TestListWithFilter` — creates 2, updates one to DOING, filter by TODO returns 1
7. `TestUpdate` — updates title, checks new title
8. `TestUpdateStatus` — updates status to DOING, checks
9. `TestDelete` — creates, deletes, Get returns error
10. `TestDeleteNotFound` — nonexistent ID returns error

**Validation gate:** `go test ./model/... -count=1` exits 0, output contains "PASS"

## Module: main.go

**Package:** `main`

**HTTP Server** on configurable port (default 8890) with REST API + embedded HTML UI.

**API Endpoints:**
```
GET  /api/tasks              → List all tasks (optional ?status=todo|doing|done filter)
POST /api/tasks              → Create task (JSON body: {title, description})
GET  /api/tasks/:id          → Get single task
PATCH /api/tasks/:id         → Update task (JSON body: {title?, description?, status?})
DELETE /api/tasks/:id        → Delete task
GET  /                       → Serve embedded HTML kanban board UI
```

**Response format:** JSON. Tasks as `{"tasks": [...]}` for list, single task object for get/create/update.
**Error format:** `{"error": "message"}` with appropriate HTTP status codes (400, 404, 405).
**CORS:** Allow all origins (`Access-Control-Allow-Origin: *`).

**Embedded HTML UI features:**
- Three columns: To Do, In Progress, Done
- Add task form (title + description)
- Move tasks between columns with buttons
- Delete button per task
- Auto-refresh every 5 seconds
- Shows task count per column

**Implementation notes:**
- Use `net/http` standard library (no framework)
- CRITICAL: Use Go backtick raw string literals (` `) for the HTML constant — NOT double quotes. Backtick strings don't interpret escape sequences, which is essential for HTML containing JavaScript.
- Global `store` variable (in-memory, no persistence)
- `handleTasks` for /api/tasks (GET list, POST create)
- `handleTask` for /api/tasks/:id (GET, PATCH, DELETE)
- `handleUI` for / (serves HTML string constant)
- HTML embedded as Go string constant `uiHTML`
- Port from environment or default 8890

**Validation gate:** `go build .` exits 0 AND `go test -count=1 -run TestUIReturnsHTML` exits 0

## Module: main_test.go

**Package:** `main` (same package — uses httptest)

**12 HTTP integration tests:**
1. `TestUIReturnsHTML` — GET / returns 200 with HTML containing "Task Board"
2. `TestCreateTaskViaAPI` — POST /api/tasks returns 201 with task JSON
3. `TestCreateTaskEmptyTitle` — POST with empty title returns 400
4. `TestListTasks` — Create 2, GET /api/tasks returns 2 tasks
5. `TestListTasksWithFilter` — Filter by status returns correct subset
6. `TestGetTask` — GET /api/tasks/:id returns correct task
7. `TestGetTaskNotFound` — GET nonexistent returns 404
8. `TestUpdateTaskStatus` — PATCH with status changes it
9. `TestUpdateTaskTitle` — PATCH with title changes it
10. `TestDeleteTask` — DELETE then GET returns 404
11. `TestDeleteTaskNotFound` — DELETE nonexistent returns 404
12. `TestEmptyListReturnsEmptyArray` — GET /api/tasks with no tasks returns {"tasks":[]}

**Validation gate:** `go test -count=1` exits 0, all 12 pass
