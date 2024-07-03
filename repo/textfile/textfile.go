package textfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
)

type RepoTextFile struct {
	fileName string
}

func (j *RepoTextFile) Add(todo model.Todo) error {
	panic("not implemented")
}

func (j *RepoTextFile) GetAll() ([]model.Todo, error) {
	panic("not implemented")
}

func (j *RepoTextFile) Get(id string) (model.Todo, error) {
	panic("not implemented")
}

func (j *RepoTextFile) GetAllStatus(status model.Status) ([]model.Todo, error) {
	panic("not implemented")
}

func (j *RepoTextFile) UpdateData(id string, newdata string) (model.Todo, error) {
	panic("not implemented")
}

func (j *RepoTextFile) UpdateStatus(id string, status model.Status) (model.Todo, error) {

	return model.Todo{}, nil
}
func (j *RepoTextFile) Remove(id string) (model.Todo, error) {
	panic("not implemented")
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

	return model.Todo{
		Id:   parts[0],
		Data: parts[1],
	}, nil
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

	s := string(b)
	s = strings.ReplaceAll(s, "\r\n", "\n")

	return linesToModel(s)
}

func modelToLine(t model.Todo) string {
	return fmt.Sprintf("%s: %s", t.Id, t.Data)
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
