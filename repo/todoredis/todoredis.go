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
func redisKeyTodo(t model.Todo) string {
	return "todo: " + t.Id
}

type RepoRedis struct {
	rd *redis.Client
}

func New(addr string, db int) repo.Repository {
	rd := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	//ctx := context.Background()

	// result, err := rd.HSet(ctx, "todo:1", "id", "", "data", "", "status", "").Result()
	// if err != nil {
	// 	panic("failed to write empty array to init file: " + err.Error())
	// }
	// _=result

	return &RepoRedis{rd: rd}
}

// Implement

// key : json
// todo: udkfdhfl " "id":"udkfdhfl","data":"yong","status":"TODO" "
func (j *RepoRedis) Add(ctx context.Context, data model.Todo) error {
	err := j.rd.HSet(ctx, redisKeyTodo(data), model.Todo{Id: data.Id, Data: data.Data, Status: data.Status}).Err()
	//err := j.rd.HSet(ctx, redisKeyTodo(data), model.Todo{}.Id, data.Id, model.Todo{}.Data, data.Data, model.Todo{}.Status, data.Status).Err()
	if err != nil {
		return fmt.Errorf("hset redis err %w", err)
	}
	return nil
}

/*
todo: {udkfdhfl} ;id "udkfdhfl",data: "yong", status: "TODO" >> keyMain = string,keyField = json

todo: {udkfdhfl}
id: "udkfdhfl"
data: "yong"
status: "TODO"
*/

func (j *RepoRedis) GetAll(ctx context.Context) ([]model.Todo, error) {
	// result2, err := j.rd.HGetAll(ctx, "todo: ").Result()
	// if err != nil {
	// 	return []model.Todo{}, fmt.Errorf("hgetall redis err %w", err)
	// }
	// _ = result2

	todos := []model.Todo{}

	keyMain, err := j.rd.Keys(ctx, "*").Result()
	if err != nil {
		return []model.Todo{}, fmt.Errorf("keys redis err %w", err)
	}

	//get = km kf > v
	//getall = km > f1,v1, f2,v2 f3,v3

	t := model.Todo{}
	for _, v := range keyMain {
		keyMainMap, err := j.rd.HGetAll(ctx, v).Result()
		if err != nil {  
			return []model.Todo{}, fmt.Errorf("hgetall redis err %w", err)
		}

		for i := range keyMainMap {
			field := keyMainMap[i]
			s, err := j.rd.HGet(ctx, field,field.).Result()
			if err != nil {
				return []model.Todo{}, fmt.Errorf("hget redis err %w", err)
			}

			t.Id = s
			t.Data = ""
		}

		// line := strings.Split(todo[value], "\n")
		// t.Id = line[1]
		// t.Data = line[3]
		// t.Status = model.Status(line[5])

		todos = append(todos, t)
	}

	return todos, nil
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
