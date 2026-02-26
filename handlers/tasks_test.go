package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"tasks-api/models"
	"tasks-api/repos"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB initializes an in-memory SQLite database for testing and returns a repos.DB instance
// It also ensures that the database connection is properly closed after the test completes.

func setupTestDB(t *testing.T) *repos.DB {
	t.Helper()

	database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	db := repos.New(database)
	if err := db.InitSchema(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Errorf("Failed to close test database: %v", err)
		}
	})
	return db
}


func TestCreate_Success (t *testing.T) {
	db := setupTestDB(t)
	handler := NewTaskHandler(db) // creates a new TaskHandler instance using the test database

	body:= `{"title": "Buy Groceries", "completed": false}`
	req:= httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	// httptest tools are used to create a response recorder and simulate an HTTP request to the Create handler
	rr := httptest.NewRecorder() // creates a fake response writer to capture the handler's output
	handler.Create(rr, req)

	// 1. Check if the status code is 201 Created
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	// 2. Check if the response contains the expected JSON with the new task ID
	var result map[string]int64
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	
	 if result["id"] != 1 {
		t.Errorf("Expected task ID 1, got %d", result["id"])
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	handler := NewTaskHandler(db) // creates a new TaskHandler instance using the test database

	body :=  `{Not a valid JSON}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder() // creates a fake response writer to capture the handler's output

	handler.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreate_EmptyBody (t *testing.T) {
	db := setupTestDB(t)
	handler := NewTaskHandler(db) // creates a new TaskHandler instance using the test database

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(""))

	rr := httptest.NewRecorder()
	handler.Create(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}


// -------------- List Tests --------------

func TestList_WithTasks(t *testing.T) {
	db := setupTestDB(t)
	handler := NewTaskHandler(db) // creates a new TaskHandler instance using the test database

	// Send the database with a couple of tasks
	db.Create("Task A", false)
	db.Create("Task B", true)
	
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()
	handler.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var tasks []models.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	titles:= map[string]bool{}

	for _, task :=range tasks {
		titles[task.Title] = true
	}

	if !titles["Task A"] || !titles["Task B"] {
		t.Errorf("Expected tasks 'Task A' and 'Task B' in response")
	}
}

func TestList_Empty(t *testing.T) {
	db := setupTestDB(t)
	handler := NewTaskHandler(db) 

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

	rr := httptest.NewRecorder()

	handler.List(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// An empty list should return null or [] - either is acceptable from json.Encode
	var tasks []models.Task
	if err :=json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
} 

// -------------- Delete Tests --------------

// TODO: Implement tests for the Delete handler, including cases for successful deletion, invalid task ID, and database errors.