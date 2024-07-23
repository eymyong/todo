package todoredis

import (
	"context"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
	"github.com/redis/go-redis/v9"
)

func redisKeyTodo(t model.Todo) string {
	return ""
}

type RepoRedis struct {
	rd *redis.Client
}

func New(addr string, db int) repo.Repository {
	rd := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx := context.Background()

	result, err := rd.HGetAll(ctx, "todo").Result()
	if err != nil || len(result) == 0 {
		err = rd.HSet(ctx, "todo:1", "id", "", "data", "", "status", "").Err()
		if err != nil {
			panic("failed to write empty array to init file: " + err.Error())
		}
	}

	return &RepoRedis{rd: rd}
}

// Implement

func (j *RepoRedis) Add(ctx context.Context, data model.Todo) error {

	panic("")
}

func (j *RepoRedis) GetAll(ctx context.Context) ([]model.Todo, error) {
	panic("")
}

func (j *RepoRedis) Get(ctx context.Context, id string) (model.Todo, error) {
	panic("")
}

func (j *RepoRedis) GetStatus(ctx context.Context, status model.Status) ([]model.Todo, error) {
	panic("")
}

func (j *RepoRedis) UpdateData(ctx context.Context, id string, newdata string) (model.Todo, error) {
	panic("")
}

func (j *RepoRedis) UpdateStatus(ctx context.Context, id string, status model.Status) (model.Todo, error) {
	panic("")
}

func (j *RepoRedis) Remove(ctx context.Context, id string) (model.Todo, error) {
	panic("")
}
