package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/eymyong/todo/cmd/api/internal/handler"
	"github.com/eymyong/todo/repo"
	"github.com/eymyong/todo/repo/jsonfile"
	"github.com/eymyong/todo/repo/jsonfilemap"
)

const JsonFile = "json"
const JsonMap = "jsonmap"
const TextFile = "text"

func initRepo() repo.Repository {
	envRepo := os.Getenv("REPO")
	envFile := os.Getenv("FILENAME")

	var repo repo.Repository
	switch envRepo {
	case JsonMap:
		if envFile == "" {
			envFile = "todo.map.json"
		}

		repo = jsonfilemap.New(envFile)

	case JsonFile:
		fallthrough

	default:
		if envFile == "" {
			envFile = "todo.json"
		}

		repo = jsonfile.New(envFile)
	}

	return repo
}

func main() {
	repo := initRepo()
	h := handler.New(repo)

	r := mux.NewRouter()
	r.HandleFunc("/get-all", h.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/get/{todo-id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/add", h.Add).Methods(http.MethodPost)
	r.HandleFunc("/delete/{todo-id}", h.Delete).Methods(http.MethodDelete)

	http.ListenAndServe(":8000", r)
}
