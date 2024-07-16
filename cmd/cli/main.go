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
/*
1f17567b-f733-4e2f-b7ef-4b4e9af6caa8: one: TODO
17aac420-1f99-4f05-844b-b77cc24f1244: two: TODO
74a398bb-7148-428b-b912-f9b4ad571bde: three: TODO
*/
type Mode string

const (
	ModeAdd          Mode = "--add"
	ModeGetAll       Mode = "--get-all"
	ModeGet          Mode = "--get"
	ModeGetStatus    Mode = "--get-status"
	ModeUpdateData   Mode = "--update"
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
			return
		}

		fmt.Printf("Get to ID: %s\nData: %s\nStatus: %s", data.Id, data.Data, data.Status)
		return

	case ModeGetStatus:
		todos, err := methodGetStatus(repo, job.status)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(todos) == 0 {
			fmt.Printf("Not found data to staus: %s", job.status)
			return
		}

		fmt.Printf("status: `%s`\n%s", job.status, todos)
		return

	case ModeUpdateData:
		old, err := methodUpdate(repo, job.id, job.data)
		if err != nil {
			fmt.Println(err)
			return

		}

		new := model.Todo{
			Id:   job.id,
			Data: job.data,
		}
		fmt.Println("Succeed")
		fmt.Printf("Old todo ID: %s, data: %s\n", old.Id, old.Data)
		fmt.Printf("New todo ID: %s, data: %s\n", new.Id, new.Data)
		return

	case ModeUpdateStatus:
		old, err := methodUpdateStatus(repo, job.id, job.status)
		if err != nil {
			fmt.Println(err)
			return
		}

		new := model.Todo{
			Id:     job.id,
			Status: job.status,
		}
		fmt.Println("Succeed")
		fmt.Printf("Old status ID: %s, status: %s\n", old.Id, old.Status)
		fmt.Printf("New status ID: %s, status: %s\n", new.Id, new.Status)
		return

	case ModeRemove:
		data, err := methodRm(repo, job.id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Succeed")
		fmt.Printf("Remove to ID: %s\ntodo: %s", data.Id, data)
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

		if args[1] == "--get-status" {
			return job{}, errors.New("there is no information to get-status")
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

		if args[1] == "--get-status" {
			return job{mode: ModeGetStatus, status: model.Status(args[2])}, nil
		}

		if args[1] == "--update" {
			return job{}, errors.New("Not data to update-data")
		}

		if args[1] == "--update-status" {
			return job{}, errors.New("Not data to update-status")
		}

		if args[1] == "--rm" {
			return job{mode: ModeRemove, id: args[2]}, nil
		}

	}

	if len(args) == 4 {
		if args[1] == "--update" {
			return job{mode: ModeUpdateData, id: args[2], data: args[3]}, nil
		}

		if args[1] == "--update-status" {
			return job{mode: ModeUpdateStatus, id: args[2], status: model.Status(args[3])}, nil
		}
	}

	return job{}, errors.New("input incorrect")

}

func methodAdd(r repo.Repository, data string) error {

	err := r.Add(model.Todo{
		Id:     uuid.NewString(),
		Data:   data,
		Status: model.StatusTodo,
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
	todo, err := r.GetStatus(status)
	if err != nil {
		return []model.Todo{}, err
	}

	return todo, nil
}

func methodUpdate(r repo.Repository, id string, newdata string) (model.Todo, error) {
	todo, err := r.UpdateData(id, newdata)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func methodUpdateStatus(r repo.Repository, id string, status model.Status) (model.Todo, error) {
	todo, err := r.UpdateStatus(id, status)
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
