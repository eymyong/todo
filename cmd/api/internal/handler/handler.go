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

type HandlerTodo struct {
	repo repo.Repository
}

func New(repo repo.Repository) *HandlerTodo {
	return &HandlerTodo{repo: repo}
}

func sendJson(w http.ResponseWriter, status int, data interface{}) { //
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func readBody(r *http.Request) ([]byte, error) { //
	defer r.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *HandlerTodo) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)   //
	id := vars["todo-id"] //
	if id == "" {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"error": "missing id",
		})

		return
	}

	todo, err := h.repo.Get(id)
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  fmt.Sprintf("failed to get todo %s", id),
			"reason": err.Error(),
		})

		return
	}

	sendJson(w, http.StatusOK, todo)
}

func (h *HandlerTodo) GetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.GetAll()
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

func (h *HandlerTodo) Add(w http.ResponseWriter, r *http.Request) {
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

	err = h.repo.Add(todo)
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  "failed to create todo",
			"reason": err.Error(),
		})

		return
	}

	sendJson(w, http.StatusCreated, map[string]interface{}{
		"success": "ok",
		"created": todo,
	})
}

func (h *HandlerTodo) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["todo-id"]
	if !ok {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"error": "missing id",
		})
		return
	}

	todo, err := h.repo.Remove(id)
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

func (h *HandlerTodo) UpdateId(w http.ResponseWriter, r *http.Request) {
	b, err := readBody(r)
	if err != nil {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"error":  "failed to read body",
			"reason": err.Error(),
		})
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["todo-id"]
	if !ok {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"err": "missing id",
		})
		return
	}

	todo, err := h.repo.Update(id, string(b))
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  fmt.Sprintf("failed to update id %s", id),
			"reason": err.Error(),
		})
		return
	}

	sendJson(w, http.StatusOK, map[string]interface{}{
		"sucess": fmt.Sprintf("update to id %s", id),
		"update": todo,
	})
}
