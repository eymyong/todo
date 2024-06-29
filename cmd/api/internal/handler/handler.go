package handler

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
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"error": "missing id",
		})

		return
	}

	todo, err := s.repo.Get(id)
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  fmt.Sprintf("failed to get todo %s", id),
			"reason": err.Error(),
		})

		return
	}

	sendJson(w, http.StatusOK, todo)
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := s.repo.GetAll()
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  "failed to get all todos",
			"reason": err.Error(),
		})

		return
	}

	m := make(map[string]string)
	for i := range todos {
		t := todos[i]
		m[t.Id] = t.Data
	}

	sendJson(w, http.StatusOK, m)
}

func (s *Server) Add(w http.ResponseWriter, r *http.Request) {
	b, err := readBody(r)
	if err != nil {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"error":  "failed to read body",
			"reason": err.Error(),
		})

		return
	}

	body := string(b)
	todo := model.Todo{
		Id:   uuid.NewString(),
		Data: body,
	}

	err = s.repo.Add(todo)
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  "failed to create todo",
			"reason": err.Error(),
		})

		return
	}

	sendJson(w, 201, map[string]interface{}{
		"success": "ok",
		"created": todo,
	})
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["todo-id"]
	if !ok {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"error": "missing id",
		})
		return
	}

	todo, err := s.repo.Remove(id)
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  fmt.Sprintf("failed to remove id %s", id),
			"reason": err.Error(),
		})
		return
	}

	sendJson(w, http.StatusOK, map[string]interface{}{
		"sucess":  "ok",
		"deleted": todo,
	})
}
