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
