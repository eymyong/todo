package main

import (
	"context"
	"fmt"

	"github.com/eymyong/todo/model"
	"github.com/redis/go-redis/v9"
)

type rRedis struct {
	rd *redis.Client
}

func main() {

	rd := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	ctx := context.Background()

	err := rd.HSet(ctx, "test:1", "id", "2", "data", "two", "status", "DONE").Err()
	if err != nil {
		panic(err)
	}

	keyStr, err := rd.Keys(ctx, "*").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("keys in redis")
	for _, v := range keyStr {
		fmt.Println(v)
	}

	todos := []model.Todo{}
	for _, v := range keyStr {
		data, err := rd.HGetAll(ctx, v).Result()
		if err != nil {
			panic(err)
		}

		fmt.Println(data)

		//datab, err := json.Marshal(data)
		//json.Unmarshal(datab)
		todo := model.Todo{}
		for k, v := range data {
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

	fmt.Println(todos)

	err = rd.Del(ctx, "test:1").Err()
	if err != nil {
		panic(err)
	}

}
