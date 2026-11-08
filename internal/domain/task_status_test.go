package domain

import (
	"encoding/json"
	"testing"
)

type CreateTaskRequest struct {
	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`
}

func TestCreateTaskRequestWithValidStatus(t *testing.T) {
	body := []byte(`{"title": "Mi tarea", "status": "todo"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)
	if err != nil {
		t.Fatalf("error not expected, but received: %v", err)
	}

	if req.Status != TaskStatusTodo {
		t.Fatalf("Expected status todo, received: %s", req.Status)
	}
}

func TestCreateTaskRequestWithInvalidStatus(t *testing.T) {
	body := []byte(`{"title": "Mi tarea", "status": "yuca"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)
	if err == nil {
		t.Fatalf("error expected, but received nil")
	}

	expected := `task status invalid: "yuca"`

	if err.Error() != expected {
		t.Fatalf("expected error %q, received instead %q", expected, err.Error())
	}
}
