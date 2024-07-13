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
	todosExpected := makeTodos()
	lines := modelToLines(todosExpected)

	err := os.WriteFile(fileName, []byte(lines), 0o664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
	}

	todosActual, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if len(todosActual) != len(todosExpected) {
		t.Errorf("unexpected length '%d', expecting '%d'", len(todosActual), len(todosExpected))
	}

	for i := range todosActual {
		actual := todosActual[i]
		expected := todosExpected[i]

		if actual != expected {
			t.Errorf("unexpected value, expecting='%+v', got='%+v'", expected, actual)
		}
	}

	err = os.WriteFile(fileName, []byte{}, 0o664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}
}

func TestReadDecode_err(t *testing.T) {
	expectederrReadFile := "no data"
	_, err := readDecode(fileNameErr)
	if err == nil {
		t.Errorf("expected err but got nil")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectederrReadFile) {
		t.Errorf("expected '%s' but got '%s'", expectederrReadFile, err.Error())
		return
	}
}

func TestModelToLine_Happy(t *testing.T) {
	todo := model.Todo{
		Id:     "1",
		Data:   "one",
		Status: model.StatusTodo,
	}
	// 1: one: TODO
	todoStr := modelToLine(todo)

	parts := strings.Split(todoStr, ": ")

	if len(parts) != 3 {
		t.Errorf("unexpected parts, expecting 3 parts, got %d", len(parts))
	}

	if parts[0] != todo.Id {
		t.Errorf("expected id: `%s` but got `%s`", todo.Id, parts[0])
	}

	if parts[1] != todo.Data {
		t.Errorf("expected data: `%s` but got `%s`", todo.Data, parts[1])
	}

	if parts[2] != string(todo.Status) {
		t.Errorf("expected status: `%s` but got `%s`", todo.Status, parts[2])
	}
}

func TestModelToLines_Happy(t *testing.T) {
	todos := []model.Todo{
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

	expected := "1: one\n2: two"
	actual := modelToLines(todos)

	if actual != expected {
		t.Errorf("unexpected value: expecting `%s` but got `%s`", expected, actual)
	}
}

func TestLineToModel_Happy(t *testing.T) {
	line := "24: foo"
	expected := model.Todo{
		Id:   "24",
		Data: "foo",
	}

	actual, err := lineToModel(line)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if actual != expected {
		t.Errorf("unexpected value, expecting='%+v', got='%+v'", expected, actual)
	}
}

func TestLinesToModel_Happy(t *testing.T) {
	lines := "1: one: TODO\n2: two: DONE\n3: three: TODO"
	expecteds := []model.Todo{
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
		{
			Id:     "3",
			Data:   "three",
			Status: model.StatusTodo,
		},
	}

	actuals, err := linesToModel(lines)
	if err != nil {
		t.Errorf("unexpexted err: `%s`", err)
		return
	}

	if len(actuals) != len(expecteds) {
		t.Errorf("unexpected length '%d', expecting '%d'", len(actuals), len(expecteds))
	}

	for i := range expecteds {
		actual := actuals[i]
		expected := expecteds[i]

		if actual != expected {
			t.Errorf("unexpected value, expecting='%+v', got='%+v'", expected, actual)
		}
	}
}
