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

func TestReadDecode_fileErr(t *testing.T) {
	expectedErrReadFile := "failed to readfile"
	_, err := readDecode(fileNameErr)
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErrReadFile) {
		t.Errorf("expected '%s' but got '%s'", expectedErrReadFile, err.Error())
		return
	}
}

// case readDecode return []model.Todo{}
func TestReadDecode_linesToModelErr(t *testing.T) {
	expectedErr := "not data"

	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Error("unexpected error", err)
	}

	_, err = linesToModel(string(b))
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
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

	expected := "1: one: TODO\n2: two: DONE"
	actual := modelToLines(todos)

	if actual != expected {
		t.Errorf("unexpected value: expecting `%s` but got `%s`", expected, actual)
	}
}

func TestLineToModel_Happy(t *testing.T) {
	line := "24: foo: DONE"
	expected := model.Todo{
		Id:     "24",
		Data:   "foo",
		Status: "DONE",
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

func TestAdd_Happy(t *testing.T) {
	data := model.Todo{
		Id:     "10",
		Data:   "yong",
		Status: model.StatusTodo,
	}

	expectedTodo := `10: yong: TODO`

	repo := RepoTextFile{
		fileName: fileName,
	}

	err := repo.Add(data)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
	}

	b, err := os.ReadFile(repo.fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
	}

	if expectedTodo != string(b) {
		t.Errorf("expected todo: `%s` but got `%s`", expectedTodo, string(b))
	}

	err = os.WriteFile(repo.fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
	}

}

func TestAdd_lineToModelErr(t *testing.T) {
	expectedErr := "not data"
	data := ""

	_, err := lineToModel(data)
	if err == nil {
		t.Errorf("expected err but fot nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
		return
	}

}

func TestGetAll_Happy(t *testing.T) {
	repo := RepoTextFile{
		fileName: fileName,
	}
	expectedTodos := []model.Todo{
		{
			Id:     "1",
			Data:   "one",
			Status: model.StatusTodo,
		},
		{
			Id:     "2",
			Data:   "two",
			Status: model.StatusTodo,
		},
		{
			Id:     "3",
			Data:   "three",
			Status: model.StatusDone,
		},
	}

	lines := modelToLines(expectedTodos)

	err := os.WriteFile(repo.fileName, []byte(lines), 0664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	actualTodos, err := repo.GetAll()
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if len(actualTodos) != len(expectedTodos) {
		t.Errorf("unexpected length '%d', expecting '%d'", len(actualTodos), len(expectedTodos))
	}

	for i := range expectedTodos {
		expected := expectedTodos[i]
		actual := actualTodos[i]

		if actual != expected {
			t.Errorf("unexpected value, expecting='%+v', got='%+v'", expected, actual)
		}
	}

	err = os.WriteFile(repo.fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

}

func TestGetAll_lineToModelErr(t *testing.T) {
	expectedErr := "not data"
	data := ""

	_, err := lineToModel(data)
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
		return
	}

}

func TestGetAll_Err(t *testing.T) {
	repo := RepoTextFile{
		fileName: fileName,
	}
	expectedErr := "not found data to file"

	_, err := repo.GetAll()
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
		return
	}

}

func TestGet_Happy(t *testing.T) {
	lines := `1: one: DONE
	2: eiei: TODO
	3: hahaha
	5: asdfk:`

	expectedTodo := model.Todo{
		Id:     "1",
		Data:   "one",
		Status: model.StatusDone,
	}

	repo := RepoTextFile{
		fileName: fileName,
	}

	err := os.WriteFile(repo.fileName, []byte(lines), 0664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	actual, err := repo.Get(expectedTodo.Id)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if expectedTodo != actual {
		t.Errorf("unexpected value, expecting='%+v', got='%+v'", expectedTodo, actual)
		return
	}

	err = os.WriteFile(repo.fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

}

func TestGet_IdErr(t *testing.T) {
	lines := `1: one: DONE
	2: eiei: TODO
	3: hahaha
	5: asdfk:`

	err := os.WriteFile(fileName, []byte(lines), 0664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	repo := RepoTextFile{
		fileName: fileName,
	}

	expectedErr := "not found id"
	expectedId := "kuy"

	_, err = repo.Get(expectedId)
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
		return
	}

}

func TestGetStatus_Happy(t *testing.T) {
	lines := `1: one: DONE
2: eiei: TODO
3: hahaha
5: asdfk:`

	expectedsTODO := []model.Todo{
		{
			Id:     "2",
			Data:   "eiei",
			Status: model.StatusTodo,
		},
		{
			Id:     "3",
			Data:   "hahaha",
			Status: model.StatusTodo,
		},
		{
			Id:     "5",
			Data:   "asdfk",
			Status: model.StatusTodo,
		},
	}

	err := os.WriteFile(fileName, []byte(lines), 0o664)
	if err != nil {
		t.Error("unexpected error", err)
	}

	repo := RepoTextFile{fileName: fileName}
	actuals, err := repo.GetStatus(model.StatusTodo)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if len(actuals) != len(expectedsTODO) {
		t.Errorf("unexpected length, expecting=%d, actual=%d", len(actuals), len(expectedsTODO))
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Error("unexpected error:", err)
		return
	}
}

func TestGet_statusErr(t *testing.T) {
	data := []model.Todo{
		{
			Id:     "2",
			Data:   "eiei",
			Status: model.StatusTodo,
		},
		{
			Id:     "3",
			Data:   "hahaha",
			Status: model.StatusTodo,
		},
		{
			Id:     "5",
			Data:   "asdfk",
			Status: model.StatusTodo,
		},
	}

	lines := modelToLines(data)

	err := os.WriteFile(fileName, []byte(lines), 0664)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	repo := RepoTextFile{
		fileName: fileName,
	}

	expectedErr := "status is not correct"
	expectedStatus := model.Status("kuyy")

	_, err = repo.GetStatus(expectedStatus)
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
		return
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Error("unexpected error", err)
		return
	}

}

func TestUpdateData_Happy(t *testing.T) {
	data := []model.Todo{
		{
			Id:     "2",
			Data:   "eiei",
			Status: model.StatusTodo,
		},
		{
			Id:     "3",
			Data:   "hahaha",
			Status: model.StatusTodo,
		},
		{
			Id:     "5",
			Data:   "asdfk",
			Status: model.StatusTodo,
		},
	}

	lines := modelToLines(data)

	err := os.WriteFile(fileName, []byte(lines), 0664)
	if err != nil {
		t.Error("unexpected error", err)
	}

	expectedTodo := model.Todo{
		Id:     "2",
		Data:   "eiei",
		Status: model.StatusTodo,
	}

	newData := "two"

	repo := RepoTextFile{fileName: fileName}
	actuals, err := repo.UpdateData(expectedTodo.Id, newData)
	if err != nil {
		t.Error("unexpected error", err)
		return
	}

	if expectedTodo.Data != actuals.Data {
		t.Errorf("expected data: `%s` but got `%s`", newData, actuals)
		return
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Error("unexpected error", err)
		return
	}

}

func TestUpdateData_IdErr(t *testing.T) {
	data := []model.Todo{
		{
			Id:     "2",
			Data:   "eiei",
			Status: model.StatusTodo,
		},
		{
			Id:     "3",
			Data:   "hahaha",
			Status: model.StatusTodo,
		},
		{
			Id:     "5",
			Data:   "asdfk",
			Status: model.StatusTodo,
		},
	}

	lines := modelToLines(data)

	err := os.WriteFile(fileName, []byte(lines), 0664)
	if err != nil {
		t.Error("unexpected error", err)
	}

	expectedErr := "not found id"
	newData := "two"

	repo := RepoTextFile{fileName: fileName}

	_, err = repo.UpdateData("10", newData)
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
		return
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Error("unexpected error", err)
	}

}

func TestUpdateStatus_Happy(t *testing.T) {
	data := []model.Todo{
		{
			Id:     "2",
			Data:   "eiei",
			Status: model.StatusTodo,
		},
		{
			Id:     "3",
			Data:   "hahaha",
			Status: model.StatusTodo,
		},
		{
			Id:     "5",
			Data:   "asdfk",
			Status: model.StatusTodo,
		},
	}

	lines := modelToLines(data)

	err := os.WriteFile(fileName, []byte(lines), 0664)
	if err != nil {
		t.Error("unexpected error", err)
	}

	newStatus := model.StatusDone
	repo := RepoTextFile{fileName: fileName}
	_, err = repo.UpdateStatus(data[0].Id, newStatus)
	if err != nil {
		t.Error("unexpected error", err)
		return
	}

	actuals, err := readDecode(fileName)
	if err != nil {
		t.Error("unexpected error", err)
		return
	}

	if newStatus != actuals[0].Status {
		t.Errorf("expect status `%s` but got `%s`", newStatus, actuals[0].Status)
		return
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Error("unexpected error", err)
	}

}

func TestRemove_Happy(t *testing.T) {
	data := []model.Todo{
		{
			Id:     "2",
			Data:   "eiei",
			Status: model.StatusTodo,
		},
		{
			Id:     "3",
			Data:   "hahaha",
			Status: model.StatusTodo,
		},
		{
			Id:     "5",
			Data:   "asdfk",
			Status: model.StatusTodo,
		},
	}

	lines := modelToLines(data)

	err := os.WriteFile(fileName, []byte(lines), 0664)
	if err != nil {
		t.Error("unexpected error", err)
	}

	idToRemove := data[0].Id
	repo := RepoTextFile{fileName: fileName}
	_, err = repo.Remove(idToRemove)
	if err != nil {
		t.Errorf("unexpectErr")
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err `%s`", err)
	}

	for _, v := range todos {
		if v.Id == idToRemove {
			t.Errorf("unexpected id: `%s`", idToRemove)
		}
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Error("unexpected error", err)
	}

}
