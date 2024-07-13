package textfile

import (
	"os"
	"strings"
	"testing"

	"github.com/eymyong/todo/model"
)

// expectedTodo := model.Todo{
// 	Id:     "1",
// 	Data:   "one",
// 	Status: model.StatusTodo,
// }
// _ = expectedTodo

const (
	fileName    = "mock/test_foo.text"
	fileNameErr = "error.text"
)

func TestReadDecode_Happy(t *testing.T) {
	expectedTodo := makeTodos()
	todoStr := modelToLines(expectedTodo)

	err := os.WriteFile(fileName, []byte(todoStr), 0o664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
	}

	todosList, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if expectedTodo[0].Id != todosList[0].Id {
		t.Errorf("expected id: `%s` but got `%s`", expectedTodo[0].Id, todosList[0].Id)
	}

	if expectedTodo[0].Data != todosList[0].Data {
		t.Errorf("expected data: `%s` but got `%s`", expectedTodo[0].Data, todosList[0].Data)
	}

	if string(expectedTodo[0].Status) != string(todosList[0].Status) {
		t.Errorf("expected status: `%s` but got `%s`", expectedTodo[0].Status, todosList[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, 0o664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}
}

func TestModelToLine_Happy(t *testing.T) {
	expectedTodo := model.Todo{
		Id:     "1",
		Data:   "one",
		Status: model.StatusTodo,
	}
	// 1: one: TODO
	todoStr := modelToLine(expectedTodo)

	parts := strings.Split(todoStr, ": ")

	if parts[0] != expectedTodo.Id {
		t.Errorf("expected id: `%s` but got `%s`", expectedTodo.Id, parts[0])
	}

	if parts[1] != expectedTodo.Data {
		t.Errorf("expected data: `%s` but got `%s`", expectedTodo.Data, parts[1])
	}

	if parts[2] != string(expectedTodo.Status) {
		t.Errorf("expected status: `%s` but got `%s`", expectedTodo.Status, parts[2])
	}
}

func TestModelToLines_Happy(t *testing.T) {
	expectedTodo := []model.Todo{
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

	// 1: one: TODO
	// 2: two: DONE
	todosListStr := modelToLines(expectedTodo)

	line := strings.Split(todosListStr, "/n")
	// [{1: one: TODO},{2: two: DONE}]
	// t.Log(line[0])
	// t.Log(line[1])

	words := strings.Split(line[0], ": ")
	if words[0] != expectedTodo[0].Id {
		t.Errorf("expected id: `%s` but got `%s`", expectedTodo[0].Id, words[0])
	}

	if words[1] != expectedTodo[0].Data {
		t.Errorf("expected data: `%s` but got `%s`", expectedTodo[1].Id, words[1])
	}
	//
	if words[2] != string(expectedTodo[0].Status) {
		t.Errorf("expected status: `%s` but got `%s`", expectedTodo[2].Id, words[2])
	}

	// for _, v := range line {
	// 	words := strings.Split(v, ": ")
	// 	if words[0] != expectedTodo[0].Id {
	// 		t.Errorf("expected id: `%s` but got `%s`", words[0], expectedTodo[0].Id)
	// 	}
	// 	if words[1] != expectedTodo[0].Data {
	// 		t.Errorf("expected data: `%s` but got `%s`", words[1], expectedTodo[0].Data)
	// 	}
	// 	if words[2] != string(expectedTodo[0].Status) {
	// 		t.Errorf("expected status: `%s` but got `%s`", words[2], expectedTodo[0].Status)
	// 	}
	// }
}

func TestLineToModel_Happy(t *testing.T) {
	data := makeTodos()

	expectedTodo := modelToLine(data[0])

	parts := strings.Split(expectedTodo, ": ")

	todo, err := lineToModel(expectedTodo)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if parts[0] != todo.Id {
		t.Errorf("expected id: '%s' but got '%s'", parts[0], todo.Id)
	}

	if parts[1] != todo.Data {
		t.Errorf("expected data: '%s' but got '%s'", parts[1], todo.Id)
	}

	if parts[2] != string(todo.Status) {
		t.Errorf("expected status: '%s' but got '%s'", parts[2], todo.Id)
	}
}

func TestLinesToModel_Happy(t *testing.T) {
	expectedTodo := makeTodos()
	lines := modelToLines(expectedTodo)

	todos, err := linesToModel(lines)
	if err != nil {
		t.Errorf("unexpexted err: `%s`", err)
		return
	}

	if expectedTodo[0].Id != todos[0].Id {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodo[0].Id, todos[0].Id)
	}

	if expectedTodo[0].Data != todos[0].Data {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodo[0].Data, todos[0].Data)
	}

	if string(expectedTodo[0].Status) != string(todos[0].Status) {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodo[0].Status, todos[0].Status)
	}
}
