package todoredis

import (
	"context"
	"fmt"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
	"github.com/redis/go-redis/v9"
)

/*
todo: {udkfdhfl}
id "udkfdhfl",
data: "yong",
status: "TODO"
*/
func redisKeyTodo(id string) string {
	return "todo: " + id
}

type RepoRedis struct {
	rd *redis.Client
}

func New(addr string) repo.Repository {
	rd := redis.NewClient(&redis.Options{
		Addr: addr,
		// DB:   db,
	})

	return &RepoRedis{rd: rd}
}

func (j *RepoRedis) Add(ctx context.Context, data model.Todo) error {
	err := j.rd.HSet(ctx, redisKeyTodo(data.Id), "id", data.Id, "data", data.Data, "status", data.Status).Err()

	if err != nil {
		return fmt.Errorf("hset redis err: %w", err)
	}
	return nil
}

func (j *RepoRedis) GetAll(ctx context.Context) ([]model.Todo, error) {
	todos := []model.Todo{}

	keyMain, err := j.rd.Keys(ctx, "*").Result()
	if err != nil {
		return []model.Todo{}, fmt.Errorf("keys redis err: %w", err)
	}

	for _, v := range keyMain {
		keyMainMap, err := j.rd.HGetAll(ctx, v).Result()
		if err != nil {
			return []model.Todo{}, fmt.Errorf("hgetall redis err: %w", err)
		}
		todo := model.Todo{}
		for k, v := range keyMainMap {
			switch k {
			case "id":
				todo.Id = v
			case "data":
				todo.Data = v
			case "status":
				todo.Status = model.Status(v)
			default:
			}
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (j *RepoRedis) Get(ctx context.Context, id string) (model.Todo, error) {
	mapStr, err := j.rd.HGetAll(ctx, redisKeyTodo(id)).Result()
	if err != nil {
		return model.Todo{}, err
	}

	todo := model.Todo{}
	for k, v := range mapStr {
		switch k {
		case "id":
			todo.Id = v
		case "data":
			todo.Data = v
		case "status":
			todo.Status = model.Status(v)
		default:
		}
	}

	return todo, nil
}

func (j *RepoRedis) GetByStatus(ctx context.Context, status model.Status) ([]model.Todo, error) {
	all, err := j.GetAll(ctx)
	if err != nil {
		return []model.Todo{}, err
	}

	targets := []model.Todo{}
	for _, v := range all {
		if v.Status == status {
			targets = append(targets, v)
		}
	}

	return targets, nil
}

func (j *RepoRedis) UpdateData(ctx context.Context, id string, newdata string) (model.Todo, error) {
	todos, err := j.GetAll(ctx)
	if err != nil {
		return model.Todo{}, err
	}

	old := model.Todo{}
	for _, v := range todos {
		if v.Id == id {
			old = v
			v.Data = newdata

			err := j.rd.HSet(ctx, redisKeyTodo(id), "data", v.Data).Err()
			if err != nil {
				return model.Todo{}, fmt.Errorf("hset redis err: %w", err)
			}

			return old, nil
		}
	}

	return model.Todo{}, fmt.Errorf("not found id: %s", id)
}

func (j *RepoRedis) UpdateStatus(ctx context.Context, id string, status model.Status) (model.Todo, error) {
	todos, err := j.GetAll(ctx)
	if err != nil {
		return model.Todo{}, err
	}

	statusOk := status.IsValid()
	if statusOk != true {
		return model.Todo{}, fmt.Errorf("bad status: %s", status)
	}

	old := model.Todo{}
	for _, v := range todos {
		if v.Id == id {
			old = v
			v.Status = status

			err := j.rd.HSet(ctx, redisKeyTodo(id), "status", v.Status).Err()
			if err != nil {
				return model.Todo{}, fmt.Errorf("hset redis err: %w", err)
			}

			return old, nil
		}
	}

	return model.Todo{}, fmt.Errorf("not found id: %s", id)
}

func (j *RepoRedis) Remove(ctx context.Context, id string) (model.Todo, error) {
	todos, err := j.GetAll(ctx)
	if err != nil {
		return model.Todo{}, err
	}

	old := model.Todo{}
	var idOk bool
	for _, v := range todos {
		if v.Id == id {
			idOk = true
			old = v
		}
	}

	if idOk == false {
		return model.Todo{}, fmt.Errorf("not found id: %s", id)
	}

	err = j.rd.Del(ctx, redisKeyTodo(id)).Err()
	if err != nil {
		return model.Todo{}, fmt.Errorf("del redis err: %w", err)
	}

	return old, nil
}
