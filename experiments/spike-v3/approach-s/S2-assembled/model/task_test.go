package model

import "testing"

func TestCreate(t *testing.T) {
	s := NewStore()
	task, err := s.Create("Test Task", "Description")
	if err != nil {
		t.Fatal(err)
	}
	if task.ID == "" {
		t.Error("expected non-empty ID")
	}
	if task.Title != "Test Task" {
		t.Errorf("expected title 'Test Task', got %q", task.Title)
	}
	if task.Status != StatusTodo {
		t.Errorf("expected status TODO, got %q", task.Status)
	}
}

func TestCreateEmptyTitle(t *testing.T) {
	s := NewStore()
	_, err := s.Create("", "")
	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestGet(t *testing.T) {
	s := NewStore()
	created, _ := s.Create("Task", "")
	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Task" {
		t.Errorf("expected 'Task', got %q", got.Title)
	}
}

func TestGetNotFound(t *testing.T) {
	s := NewStore()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestList(t *testing.T) {
	s := NewStore()
	s.Create("Task 1", "")
	s.Create("Task 2", "")
	tasks := s.List(nil)
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestListWithFilter(t *testing.T) {
	s := NewStore()
	s.Create("Task 1", "")
	task2, _ := s.Create("Task 2", "")
	doing := StatusDoing
	s.Update(task2.ID, nil, nil, &doing)
	
	todo := StatusTodo
	tasks := s.List(&todo)
	if len(tasks) != 1 {
		t.Errorf("expected 1 TODO task, got %d", len(tasks))
	}
}

func TestUpdate(t *testing.T) {
	s := NewStore()
	task, _ := s.Create("Original", "")
	newTitle := "Updated"
	updated, err := s.Update(task.ID, &newTitle, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Updated" {
		t.Errorf("expected 'Updated', got %q", updated.Title)
	}
}

func TestUpdateStatus(t *testing.T) {
	s := NewStore()
	task, _ := s.Create("Task", "")
	doing := StatusDoing
	updated, _ := s.Update(task.ID, nil, nil, &doing)
	if updated.Status != StatusDoing {
		t.Errorf("expected DOING, got %q", updated.Status)
	}
}

func TestDelete(t *testing.T) {
	s := NewStore()
	task, _ := s.Create("Task", "")
	err := s.Delete(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Get(task.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := NewStore()
	err := s.Delete("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}
