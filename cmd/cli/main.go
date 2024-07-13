package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
	"github.com/eymyong/todo/repo/jsonfile"
	"github.com/eymyong/todo/repo/jsonfilemap"
	"github.com/eymyong/todo/repo/textfile"
	"github.com/google/uuid"
)

/*
git remote origin https://github.com/eymyong/todo.git
git add -A
git commit -m "homework"
git push origin

*/

/*
[main.exe --add "kuyhee"]   ==> {"id": 2, "text": "kuyhee"}
main.exe --get 2          ==> {"id": 2, "text": "kuyhee"}
main.exe --rm 2
main.exe --update  2 "heekuy"  ==> {"id": 2, "text": "heekuy"}
main.exe                  ==> []

$Env:REPO = "text"
$Env:FILENAME = ""
*/
type Mode string

const (
	ModeAdd          Mode = "--add"
	ModeGetAll       Mode = "--get-all"
	ModeGet          Mode = "--get"
	ModeGetStatus    Mode = "--get-status"
	ModeUpdate       Mode = "--update"
	ModeUpdateStatus Mode = "--update-status"
	ModeRemove       Mode = "--rm"
)

type job struct {
	id     string
	data   string
	status model.Status
	mode   Mode
}

const JsonFile = "json"
const JsonMap = "jsonmap"
const TextFile = "text"

func initRepo() repo.Repository {
	envRepo := os.Getenv("REPO")
	envFile := os.Getenv("FILENAME")

	var repo repo.Repository

	switch envRepo {
	case JsonMap:
		if envFile == "" {
			envFile = "todo.map.json"
		}

		repo = jsonfilemap.New(envFile)

	case TextFile:
		if envFile == "" {
			envFile = "todo.text"
		}

		repo = textfile.New(envFile)

	default:
		if envFile == "" {
			envFile = "todo.json"
		}

		repo = jsonfile.New(envFile)
	}

	return repo
	// switch envRepo {
	// case JsonMap:
	// 	if envFile == "" {
	// 		envFile = "todo.map.json"
	// 	}

	// 	repo = jsonfilemap.New(envFile)

	// case JsonFile:
	// 	fallthrough

	// default:
	// 	if envFile == "" {
	// 		envFile = "todo.json"
	// 	}

	// 	repo = jsonfile.New(envFile)
	// }

	// return repo

}

func main() {
	args := os.Args
	job, err := parse(args)
	if err != nil {
		panic(err)
	}

	repo := initRepo()

	switch job.mode {
	case ModeAdd:
		err = methodAdd(repo, job.data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Succeed")
		return
	case ModeGetAll:
		todoList, err := methodGetAll(repo)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(todoList) == 0 {
			fmt.Println("No data")
			return
		}

		for _, todo := range todoList {
			fmt.Println(todo)
		}

	case ModeGet:
		data, err := methodGet(repo, job.id)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Get to ID: %s\nData: %s\nStatus: %s", data.Id, data.Data, data.Status)
		return

	case ModeGetStatus:

	case ModeUpdate:
		old, err := methodUpdate(repo, job.id, job.data)
		if err != nil {
			fmt.Println(err)
		}

		new := model.Todo{
			Id:   job.id,
			Data: job.data,
		}

		fmt.Printf("Old todo ID: %s, data: %s\n", old.Id, old.Data)
		fmt.Printf("New todo ID: %s, data: %s\n", new.Id, new.Data)
		return

	case ModeUpdateStatus:

	case ModeRemove:
		data, err := methodRm(repo, job.id)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Remove to ID: %s\nData: %s", data.Id, data.Data)
		return
	default:
		fmt.Println("Incorrect Mode")
		return
	}

}

func parse(args []string) (job, error) {
	if len(args) == 1 {
		return job{mode: ModeGetAll}, nil
	}

	if len(args) == 2 {
		if args[1] == "--add" {
			return job{}, errors.New("there is no information to add")
		}
		if args[1] == "--get" {
			return job{}, errors.New("there is no information to get")
		}
		if args[1] == "--update" {
			return job{}, errors.New("there is no information to update")
		}
		if args[1] == "--rm" {
			return job{}, errors.New("there is no information to rm")
		}
	}

	if len(args) == 3 {
		if args[1] == "--add" {
			return job{mode: ModeAdd, data: args[2]}, nil
		}
		if args[1] == "--get" {
			return job{mode: ModeGet, id: args[2]}, nil
		}
		//
		if args[1] == "get-status" {
			return job{mode: ModeGetStatus, status: model.Status(args[2])}, nil
		}
		if args[1] == "--rm" {
			return job{mode: ModeRemove, id: args[2]}, nil
		}
		if args[1] == "--update" {
			fmt.Println("Not data to update")
			return job{}, nil
		}

	}

	if len(args) == 4 {
		// if args[1] == "--update-status" {
		// 	return job{mode: ModeUpdate, id: args[2], data: args[3]}, nil
		// }
		if args[1] == "--update" {
			return job{mode: ModeUpdate, id: args[2], data: args[3]}, nil
		}

	}

	return job{}, errors.New("input incorrect")

}

func methodUpdate(r repo.Repository, id string, newdata string) (model.Todo, error) {
	todo, err := r.UpdateData(id, newdata)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func methodRm(r repo.Repository, id string) (model.Todo, error) {
	todo, err := r.Remove(id)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func methodAdd(r repo.Repository, data string) error {

	err := r.Add(model.Todo{
		Id:   uuid.NewString(),
		Data: data,
	})
	if err != nil {
		return err
	}
	return nil
}

func methodGetAll(r repo.Repository) ([]model.Todo, error) {
	todoList, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	return todoList, nil
}

func methodGet(r repo.Repository, id string) (model.Todo, error) {
	todo, err := r.Get(id)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func methodGetStatus(r repo.Repository, status model.Status) ([]model.Todo, error) {
	todo, err := r.GetStatus(model.StatusDone)
	if err != nil {
		return []model.Todo{}, err
	}

	return todo, nil
}
