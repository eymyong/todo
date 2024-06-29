package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
)

type Server struct {
	repo repo.Repository
}

func New(repo repo.Repository) *Server {
	return &Server{repo: repo}
}

func sendJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Server) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["todo-id"]
	if id == "" {
		sendJson(w, 400, map[string]interface{}{
			"error": "missing id",
		})
		return
	}

	todo, err := s.repo.Get(id)
	if err != nil {
		sendJson(w, 500, map[string]interface{}{
			"error":  fmt.Sprintf("failed to get todo %s", id),
			"reason": err.Error(),
		})
		return
	}

	sendJson(w, 200, todo)
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := s.repo.GetAll()
	if err != nil {
		errMsg := fmt.Sprintf("error reading from repo: %s", err.Error())
		sendJson(w, 500, errMsg)
		return
	}

	m := make(map[string]string)
	for i := range todos {
		t := todos[i]
		m[t.Id] = t.Data
	}

	sendJson(w, 200, m)
}

func (s *Server) Add(w http.ResponseWriter, r *http.Request) {
	b, err := readBody(r)
	if err != nil {
		errMsg := fmt.Sprintf("error readbody to json: %s", err.Error())
		sendJson(w, 500, errMsg)
		return
	}

	body := string(b)
	todo := model.Todo{
		Id:   uuid.NewString(),
		Data: body,
	}

	err = s.repo.Add(todo)
	if err != nil {
		errMsg := fmt.Sprintf("error add repo: %s", err.Error())
		sendJson(w, 500, errMsg)
		return
	}

	sendJson(w, 200, todo)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["todo-id"]
	if !ok {
		sendJson(w, 400, "missing id")
		return
	}

	todo, err := s.repo.Remove(id)
	if err != nil {
		errMsg := fmt.Sprintf("failed to remove id %s: %s\n", id, err.Error())
		sendJson(w, 500, errMsg)
		return
	}

	sendJson(w, 200, map[string]interface{}{
		"sucess":  "ok",
		"deleted": todo,
	})
}
