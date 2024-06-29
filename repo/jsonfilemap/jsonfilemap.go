package jsonfilemap

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
)

type RepoJsonFileMap struct {
	fileName string
}

func readDecode(fname string) (map[string]model.Todo, error) {
	j, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read jsonfile: %w", err)
	}

	todos := make(map[string]model.Todo)
	err = json.Unmarshal(j, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (j *RepoJsonFileMap) Add(todo model.Todo) error {

	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return err
	}

	todoMap[todo.Id] = todo

	todoByte, err := json.Marshal(todoMap)
	if err != nil {
		return err
	}

	err = os.WriteFile(j.fileName, todoByte, 0664)
	if err != nil {
		return fmt.Errorf("fail to writefile %s: %w", j.fileName, err)
	}

	return nil
}

func (j *RepoJsonFileMap) GetAll() ([]model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return []model.Todo{}, err
	}

	todoList := []model.Todo{}
	for _, todo := range todoMap {
		todoList = append(todoList, todo)
	}

	return todoList, nil
}

func (j *RepoJsonFileMap) Get(id string) (model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	todo, ok := todoMap[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("no id: %s", id)
	}

	return todo, nil
}

func (j *RepoJsonFileMap) Update(id string, newdata string) (model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	old, ok := todoMap[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("no id: %s", id)
	}

	todoMap[id] = model.Todo{Id: id, Data: newdata}
	todoMapByte, err := json.Marshal(todoMap)
	if err != nil {
		return model.Todo{}, err
	}

	err = os.WriteFile(j.fileName, todoMapByte, 0664)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to writefile jsonfile: %w", err)
	}

	return old, nil
}

func (j *RepoJsonFileMap) Remove(id string) (model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	todo, ok := todoMap[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("no id: %s", id)
	}
	delete(todoMap, id)

	todoMapBytes, err := json.Marshal(todoMap)
	if err != nil {
		return model.Todo{}, err
	}

	err = os.WriteFile(j.fileName, todoMapBytes, 0664)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func New(fileName string) repo.Repository {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil || len(fileBytes) == 0 {
		err := os.WriteFile(fileName, []byte("{}"), os.ModePerm)
		if err != nil {
			panic("failed to write empty array to init file: " + err.Error())
		}
	}

	return &RepoJsonFileMap{
		fileName: fileName,
	}
}
