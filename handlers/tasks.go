package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tasks-api/repos"
)

type TaskHandler struct {
	db *repos.DB
}

func NewTaskHandler(db *repos.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.db.List()

	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := h.db.Create(input.Title, input.Completed)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.db.Delete(id); err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
