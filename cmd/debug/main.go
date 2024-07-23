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

	err := rd.HSet(ctx, "person:1", "name", "yong", "age", "10").Err()
	// err := rd.SetEx(ctx, "testgo", "newdata2", time.Second*5).Err()
	if err != nil {
		panic(err)
	}

	// result, err := rd.HGet(ctx, "person:1", "name").Result()
	// if err != nil {
	// 	panic(err)
	// }

	//repo := rRedis{rd : rd}
	repo := repoRedis{rd: rd}
	err = repo.addTodo(ctx, model.Todo{
		Id:     "2222",
		Data:   "twotwotwo",
		Status: "DONE",
	})
	if err != nil {
		panic(err)
	}

	status, err := repo.getStatus(ctx, "2222")
	if err != nil {
		panic(err)
	}

	fmt.Println("status", status)

	// result, err := rd.Get(ctx, "zz").Result()
	// if err != nil {
	// 	if err != redis.Nil {
	// 		panic(err)
	// 	}
	// }

	// fmt.Println(result)
}
