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
	w.Header().Set("Content-Type", "application/json") //
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

	todo, err := h.repo.Get(nil, id)
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
	todos, err := h.repo.GetAll(nil)
	if err != nil {
		sendJson(w, http.StatusInternalServerError, map[string]interface{}{
			"error":  "failed to get all todos",
			"reason": err.Error(),
		})

		return
	}

	// m := make(map[string]string)
	// for i := range todos {
	// 	t := todos[i]
	// 	m[t.Id] = t.Data
	// }

	sendJson(w, http.StatusOK, todos)
}

func (h *HandlerTodo) GetAllStatus(w http.ResponseWriter, r *http.Request) {
	b, err := readBody(r)
	if err != nil {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"error":  "failed to read body",
			"reason": err.Error(),
		})
	}

	type req struct {
		Status model.Status `json:"status"`
	}

	var rr req
	err = json.Unmarshal(b, &rr)
	if err != nil {
		sendJson(w, 400, map[string]interface{}{
			"err":    "unmarshal body error",
			"reason": err.Error(),
		})
		return
	}

	if !rr.Status.IsValid() {
		sendJson(w, 400, map[string]interface{}{
			"err":    "invalid status",
			"reason": fmt.Sprintf("bad status '%s'", rr.Status),
		})
		return
	}

	if rr.Status == "" {
		rr.Status = model.StatusTodo
	}

	statusTodoList, err := h.repo.GetStatus(nil, rr.Status)
	if err != nil {
		sendJson(w, 400, map[string]interface{}{
			"err":    "invalid status",
			"reason": err.Error(),
		})
		return
	}

	sendJson(w, http.StatusOK, statusTodoList)

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
		Id:     uuid.NewString(),
		Data:   body,
		Status: model.StatusTodo,
	}

	err = h.repo.Add(nil, todo)
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

	todo, err := h.repo.Remove(nil, id)
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
	id := vars["todo-id"]
	if id == "" {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"err": "missing id",
		})
		return
	}

	todo, err := h.repo.UpdateData(nil, id, string(b))
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

// {"status":"done"}
func (h *HandlerTodo) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["todo-id"]
	if id == "" {
		sendJson(w, http.StatusBadRequest, map[string]interface{}{
			"err": "missing id",
		})
		return
	}

	b, err := readBody(r)
	if err != nil {
		sendJson(w, 400, map[string]interface{}{
			"err":    "read body error",
			"reason": err.Error(),
		})
		return
	}

	type req struct {
		Status model.Status `json:"status"`
	}

	var rr req
	err = json.Unmarshal(b, &rr)
	if err != nil {
		sendJson(w, 400, map[string]interface{}{
			"err":    "unmarshal body error",
			"reason": err.Error(),
		})
		return
	}

	if !rr.Status.IsValid() {
		sendJson(w, 400, map[string]interface{}{
			"err":    "invalid status",
			"reason": fmt.Sprintf("bad status '%s'", rr.Status),
		})
		return
	}

	if rr.Status == "" {
		rr.Status = model.StatusTodo
	}

	status, err := h.repo.UpdateStatus(nil, id, rr.Status)
	if err != nil {
		sendJson(w, 500, map[string]interface{}{
			"err":    "update-status error",
			"reason": err.Error(),
		})
		return
	}

	sendJson(w, http.StatusOK, map[string]interface{}{
		"success":       fmt.Sprintf("update-status to id '%s'", id),
		"update-status": status,
	})
}
