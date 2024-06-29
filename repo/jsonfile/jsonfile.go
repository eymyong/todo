package jsonfile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
)

type RepoJsonFile struct {
	fileName string
}

func readDecode(fname string) ([]model.Todo, error) {
	j, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read jsonfile: %w", err)
	}

	todos := []model.Todo{}
	err = json.Unmarshal(j, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (j *RepoJsonFile) Add(todo model.Todo) error {
	todoList, err := readDecode(j.fileName)
	if err != nil {
		return err
	}

	todoList = append(todoList, todo)
	out, err := json.Marshal(todoList)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(j.fileName, out, 0664)
	if err != nil {
		panic(err)
	}

	return nil
}

func (j *RepoJsonFile) GetAll() ([]model.Todo, error) {
	return readDecode(j.fileName)
}

func (j *RepoJsonFile) Get(id string) (model.Todo, error) {
	todoList, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to get jsonfile: %w", err)
	}

	for _, todo := range todoList {
		if todo.Id == id {
			return todo, nil
		}
	}

	return model.Todo{}, fmt.Errorf("no id: %s", id)
}

func (j *RepoJsonFile) Update(id string, newdata string) (model.Todo, error) {
	todoList, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to update jsonfile: %w", err)
	}

	newTodoLists := []model.Todo{}
	old := model.Todo{}
	for _, todo := range todoList {
		if id == todo.Id {
			old = todo
			todo.Data = newdata
			newTodoLists = append(newTodoLists, todo)
			continue
		}

		newTodoLists = append(newTodoLists, todo)
	}

	todoByte, err := json.Marshal(newTodoLists)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to marshal jsonfile: %w", err)
	}

	err = os.WriteFile(j.fileName, todoByte, 0664)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to writefile jsonfile: %w", err)
	}

	return old, nil
}

func (j *RepoJsonFile) Remove(id string) (model.Todo, error) {
	todoList, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to unmarshal jsonfile: %w", err)
	}

	newTodoList := []model.Todo{}
	old := model.Todo{}
	for _, todo := range todoList {
		if id == todo.Id {
			old = todo
			continue
		}
		newTodoList = append(newTodoList, todo)
	}

	todoBytes, err := json.Marshal(newTodoList)
	if err != nil {
		return model.Todo{}, err
	}

	err = os.WriteFile(j.fileName, todoBytes, 0664)
	if err != nil {
		return model.Todo{}, err
	}

	return old, nil
}

func New(fileName string) repo.Repository {
	b, err := os.ReadFile(fileName)
	if err != nil || len(b) == 0 {
		err := os.WriteFile(fileName, []byte("[]"), os.ModePerm)
		if err != nil {
			panic("failed to write empty array to init file: " + err.Error())
		}
	}
	return &RepoJsonFile{
		fileName: fileName,
	}
}
