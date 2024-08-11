package repo

import (
	"context"

	"github.com/eymyong/todo/model"
)

type Repository interface {
	Add(ctx context.Context, data model.Todo) error
	GetAll(ctx context.Context) ([]model.Todo, error)
	Get(ctx context.Context, id string) (model.Todo, error)
	GetByStatus(ctx context.Context, status model.Status) ([]model.Todo, error)
	UpdateData(ctx context.Context, id string, newdata string) (model.Todo, error)
	UpdateStatus(ctx context.Context, id string, status model.Status) (model.Todo, error)
	Remove(ctx context.Context, id string) (model.Todo, error)
}
