package repo

import "github.com/eymyong/todo/model"

type Repository interface {
	Add(model.Todo) error
	GetAll() ([]model.Todo, error)
	Get(string) (model.Todo, error)
	Update(id string, newdata string) (model.Todo, error)
	Remove(id string) (model.Todo, error)
}
