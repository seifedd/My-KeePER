package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"tasks-api/handlers"
	"tasks-api/repos"

	_ "github.com/mattn/go-sqlite3"
)

var db *repos.DB

func initDB() {
	var err error
	database, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	db = repos.New(database)
	if err := db.InitSchema(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()

	taskHandler := handlers.NewTaskHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", taskHandler.List)
	mux.HandleFunc("POST /tasks", taskHandler.Create)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.Delete)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
