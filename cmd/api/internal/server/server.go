package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eymyong/TODO-CLI/model"
	"github.com/eymyong/TODO-CLI/repo"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	repo repo.Repository
}

func New(repo repo.Repository) *Server {
	return &Server{repo: repo}
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

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := s.repo.GetAll()
	if err != nil {
		errMsg := fmt.Sprintf("error reading from repo: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))

		return
	}

	m := make(map[string]string)
	for i := range todos {
		t := todos[i]
		m[t.Id] = t.Data
	}

	b, err := json.Marshal(m)
	if err != nil {
		errMsg := fmt.Sprintf("error marshaling to json: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *Server) Add(w http.ResponseWriter, r *http.Request) {
	b, err := readBody(r)
	if err != nil {
		errMsg := fmt.Sprintf("error readbody to json: %s", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(errMsg))
		return
	}

	body := string(b)
	todo := model.Todo{
		Id:   uuid.NewString(),
		Data: body,
	}

	err = s.repo.Add(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error add repo: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["todo-id"]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("missing id"))
		return
	}

	todo, err := s.repo.Remove(id)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to remove id %s: %s\n", id, err.Error())
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sucess":  "ok",
		"deleted": todo,
	})
}
