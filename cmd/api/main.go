package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/eymyong/todo/cmd/api/internal/handler"
	"github.com/eymyong/todo/repo"
	"github.com/eymyong/todo/repo/jsonfile"
	"github.com/eymyong/todo/repo/jsonfilemap"
	"github.com/eymyong/todo/repo/textfile"
	"github.com/eymyong/todo/repo/todoredis"
)

const JsonFile = "json"
const JsonMap = "jsonmap"
const TextFile = "text"
const Redis = "redis"

func initRepo() repo.Repository {
	envRepo := os.Getenv("REPO")
	envFile := os.Getenv("FILENAME")

	var repo repo.Repository
	switch envRepo {
	case JsonFile:
		if envFile == "" {
			envFile = "todo.json"
		}
		repo = jsonfile.New(envFile)

	case JsonMap:
		if envFile == "" {
			envFile = "todo.map.json"
		}
		repo = jsonfilemap.New(envFile)

	case TextFile:
		if envFile == "" {
			envFile = "todo.text"
		}
		repo = textfile.New(envFile)

	case Redis:
		repo = todoredis.New("127.0.0.1:6379")
	}

	return repo
}

func main() {
	repo := initRepo()
	h := handler.New(repo)

	r := mux.NewRouter()
	r.HandleFunc("/get-all", h.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/get-all-status", h.GetAllStatus).Methods(http.MethodGet)
	r.HandleFunc("/get/{todo-id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/add", h.Add).Methods(http.MethodPost)
	r.HandleFunc("/delete/{todo-id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/update/{todo-id}", h.UpdateId).Methods(http.MethodPatch)
	r.HandleFunc("/update-status/{todo-id}", h.UpdateStatus).Methods(http.MethodPatch)

	http.ListenAndServe(":8000", r)
}
