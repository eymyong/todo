package main

import (
	"log"
	"net/http"
	"os"

	"github.com/eymyong/TODO-CLI/cmd/api/internal/server"
	"github.com/eymyong/TODO-CLI/repo"
	"github.com/eymyong/TODO-CLI/repo/jsonfile"
	"github.com/eymyong/TODO-CLI/repo/jsonfilemap"
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

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World\n"))
}

func yong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("kuy\n"))
}

func main() {
	repo := initRepo()
	serv := server.New(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/yong", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("kuy\n"))
	})

	mux.HandleFunc("/hello", helloWorld)
	mux.HandleFunc("/get-all", serv.GetAll)
	mux.HandleFunc("/add", serv.Add)

	// host/add

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
