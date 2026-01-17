package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tasks-api/repos"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

var db *repos.DB

func initDB() {
	var err error
	database, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	db = repos.New(database)
	if err:=db.InitSchema(); err != nil {
		log.Fatal(err)
	}
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string
		Completed bool
	}
	if err:= json.NewDecoder(r.Body).Decode(&input); err!=nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
	}

	id, err := db.Create(input.Title,input.Completed)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    if err := db.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func main() {
initDB()

initDB()
    mux := http.NewServeMux()

    mux.HandleFunc("GET /tasks", listTasks)
    mux.HandleFunc("POST /tasks", createTask)
    mux.HandleFunc("DELETE /tasks/{id}", deleteTask) 

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
