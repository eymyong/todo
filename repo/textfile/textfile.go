package textfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
)

/*
	data

1: one
2: two
3: three
*/
type RepoTextFile struct {
	fileName string
}

func (j *RepoTextFile) Add(todo model.Todo) error {
	todosList, err := readDecode(j.fileName)
	if err != nil {
		return err
	}

	todosList = append(todosList, todo)
	todosStr := modelToLines(todosList)

	err = os.WriteFile(j.fileName, []byte(todosStr), 0o664)
	if err != nil {
		return err
	}

	return nil
}

func (j *RepoTextFile) GetAll() ([]model.Todo, error) {
	todosList, err := readDecode(j.fileName)
	if err != nil {
		return []model.Todo{}, err
	}
	return todosList, nil
}

func (j *RepoTextFile) Get(id string) (model.Todo, error) {
	todosList, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	for _, v := range todosList {
		if id == v.Id {
			fmt.Println(v)
			return v, nil
		}
	}

	return model.Todo{}, nil
}

func (j *RepoTextFile) GetStatus(status model.Status) ([]model.Todo, error) {
	todosList, err := readDecode(j.fileName)
	if err != nil {
		return []model.Todo{}, err
	}

	newTodoList := []model.Todo{}
	for _, v := range todosList {
		if status == v.Status {
			newTodoList = append(newTodoList, v)
		}
	}

	return newTodoList, nil
}

func New(fileName string) repo.Repository {
	b, err := os.ReadFile(fileName)
	if err != nil || len(b) == 0 {
		err := os.WriteFile(fileName, []byte(""), os.ModePerm)
		if err != nil {
			panic("failed to write empty array to init file: " + err.Error())
		}
	}
	return &RepoTextFile{
		fileName: fileName,
	}
}

func lineToModel(line string) (model.Todo, error) {
	parts := strings.Split(line, ": ")
	if len(parts) < 2 {
		return model.Todo{}, fmt.Errorf("not data")
	}

	status := model.StatusTodo

	if len(parts) >= 3 {
		status = model.Status(parts[2])
	}

	todo := model.Todo{
		Id:     parts[0],
		Data:   parts[1],
		Status: status,
	}

	if todo.Status == "" {
		todo.Status = model.StatusTodo
	}

	return todo, nil
}

func linesToModel(data string) ([]model.Todo, error) {
	lines := strings.Split(data, "\n")

	todos := []model.Todo{}
	for _, v := range lines {
		todo, err := lineToModel(v)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func readDecode(fname string) ([]model.Todo, error) {
	b, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return []model.Todo{}, nil
	}

	s := string(b)
	s = strings.ReplaceAll(s, "\r\n", "\n")

	return linesToModel(s)
}

func modelToLine(t model.Todo) string {
	return fmt.Sprintf("%s: %s: %s", t.Id, t.Data, t.Status)
}

func modelToLines(todos []model.Todo) string {
	s := ""
	last := len(todos) - 1
	for i := range todos {
		s += modelToLine(todos[i])

		if i == last {
			continue
		}

		s += "\n"
	}

	return s
}

func modelToLinesJoin(todos []model.Todo) string {
	lines := make([]string, len(todos))
	for i := range todos {
		lines[i] = modelToLine(todos[i])
	}

	return strings.Join(lines, "\n")
}

func (j *RepoTextFile) UpdateData(id string, newdata string) (model.Todo, error) {
	return model.Todo{}, nil
}

func (j *RepoTextFile) UpdateStatus(id string, status model.Status) (model.Todo, error) {
	return model.Todo{}, nil
}

func (j *RepoTextFile) Remove(id string) (model.Todo, error) {
	return model.Todo{}, nil
}

func makeTodos() []model.Todo {
	newTodos := []model.Todo{
		{
			Id:     "1",
			Data:   "one",
			Status: model.StatusTodo,
		},
		{
			Id:     "2",
			Data:   "two",
			Status: model.StatusDone,
		},
	}
	return newTodos
}
