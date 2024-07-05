package repo

import "github.com/eymyong/todo/model"

type Repository interface {
	Add(model.Todo) error
	GetAll() ([]model.Todo, error)
	Get(string) (model.Todo, error)
	GetStatus(status model.Status) ([]model.Todo, error)
	UpdateData(id string, newdata string) (model.Todo, error)
	UpdateStatus(id string, status model.Status) (model.Todo, error)
	Remove(id string) (model.Todo, error)
}
