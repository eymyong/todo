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

// ถ้าไม่มีข้อมูลในไฟล์ return []model.Todo{}
func readDecode(fname string) ([]model.Todo, error) {
	j, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read jsonfile: %w", err)
	}

	if len(j) == 0 {
		return []model.Todo{}, nil
	}

	todos := []model.Todo{}
	err = json.Unmarshal(j, &todos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return todos, nil
}

func (j *RepoJsonFile) Add(todo model.Todo) error {
	todoList, err := readDecode(j.fileName)
	if err != nil {
		return fmt.Errorf("failed to add jsonfile: %w", err)
	}

	todoList = append(todoList, todo)
	byteTodoList, err := json.Marshal(todoList)
	if err != nil {
		return fmt.Errorf("failed to add marshal: %w", err)
	}

	err = os.WriteFile(j.fileName, byteTodoList, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to add writefile: %w", err)
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

func (j *RepoJsonFile) GetStatus(status model.Status) ([]model.Todo, error) {
	todoList, err := readDecode(j.fileName)
	if err != nil {
		return []model.Todo{}, fmt.Errorf("failed to get-status jsonfile: %w", err)
	}

	var statusTodoList []model.Todo
	for _, v := range todoList {
		if v.Status == status {
			statusTodoList = append(statusTodoList, v)
		}
	}

	return statusTodoList, nil
}

func (j *RepoJsonFile) UpdateData(id string, newdata string) (model.Todo, error) {
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
		return model.Todo{}, fmt.Errorf("failed to updatedata marshal: %w", err)
	}

	err = os.WriteFile(j.fileName, todoByte, 0664)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to updatedata writefile: %w", err)
	}

	return old, nil
}

func (j *RepoJsonFile) UpdateStatus(id string, status model.Status) (model.Todo, error) {
	todos, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	var old *model.Todo

	for i := range todos {
		t := &todos[i]
		if id == t.Id {
			old = t
			t.Status = status
		}

	}

	if old == nil {
		return model.Todo{}, fmt.Errorf("id '%s' not found", id)
	}

	todosByte, err := json.Marshal(todos)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to update-status marshal: %w", err)
	}
	err = os.WriteFile(j.fileName, todosByte, 0664)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to update-statu writefile: %w", err)
	}

	return *old, nil
}

func (j *RepoJsonFile) Remove(id string) (model.Todo, error) {
	todoList, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to Remove: %w", err)
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

	err = os.WriteFile(j.fileName, todoBytes, os.ModePerm)
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
