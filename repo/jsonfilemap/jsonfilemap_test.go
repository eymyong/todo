package jsonfilemap

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/eymyong/todo/model"
)

const fileName = "mock/test_foo.json"

func TestReadDecode_Happy(t *testing.T) {
	expectedTodosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		expectedTodosMap[key] = value
	}

	b, err := json.Marshal(expectedTodosMap)
	if err != nil {
		t.Errorf("unexpected err to marshal: `%s`", err)
		return
	}

	err = os.WriteFile(fileName, b, 0664)
	if err != nil {
		t.Errorf("unexpected err to writefile: `%s`", err)
		return
	}

	todosMap, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if len(expectedTodosMap) != len(todosMap) {
		t.Errorf("unexpected length expectedTodosMap != length todosMap")
		return
	}

	arrKeyTodosMap := []string{}
	for k := range todosMap {
		arrKeyTodosMap = append(arrKeyTodosMap, k)
	}

	for _, v := range arrKeyTodosMap {
		_, ok := expectedTodosMap[v]
		if !ok {
			t.Errorf("not found key")
			return
		}
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestReadDecodeFileName_Err(t *testing.T) {
	expectedErr := "failed to read jsonfile"

	_, err := readDecode("")
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
	}

}

func TestReadDecodeUnmarshal_Err(t *testing.T) {
	expectedErr := "unexpected end of JSON input"

	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
	}

	var u map[string]model.Todo
	err = json.Unmarshal(b, &u)
	if err == nil {
		t.Errorf("expected err but got nil")
	}

	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected err: `%s` but got `%s`", expectedErr, err.Error())
	}

}

func TestWriteEncode_Happy(t *testing.T) {
	expectedTodosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		expectedTodosMap[key] = value
	}

	if len(expectedTodosMap) != len(todos) {
		t.Errorf("expect length expectedTodosMap: `%d` but got length todos: `%d`", len(expectedTodosMap), len(todos))
		return
	}

	err := writeEncode(fileName, expectedTodosMap)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err)
		return
	}

	todosMap, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if len(expectedTodosMap) != len(todosMap) {
		t.Errorf("expect length expectedTodosMap: `%d` but got length todosMap: `%d`", len(expectedTodosMap), len(todosMap))
		return
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestAdd_Happy(t *testing.T) {
	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	newData := model.Todo{
		Id:     "10",
		Data:   "ten",
		Status: model.StatusTodo,
	}

	repo := RepoJsonFileMap{fileName: fileName}
	err = repo.Add(nil, newData)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err)
		return
	}

	newTodosMap, err := readDecode(repo.fileName)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err)
		return
	}

	_, ok := newTodosMap[newData.Id]
	if !ok {
		t.Errorf("not found key: `%s`", newData.Id)
		return
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestGetAll_Happy(t *testing.T) {
	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	repo := RepoJsonFileMap{fileName: fileName}
	newTodos, err := repo.GetAll(nil)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	for i := range todos {
		expected := todos[i]
		actual := newTodos[i]
		if expected != actual {
			t.Errorf("unexpected value, expecting='%+v', got='%+v'", expected, actual)
			return
		}
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}
}

func TestGetData_Happy(t *testing.T) {
	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	get := todos[0]

	repo := RepoJsonFileMap{fileName: fileName}
	todo, err := repo.Get(nil, get.Id)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if get != todo {
		t.Errorf("expected todo: `%+v` but got %+v", get, todo)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestGetStatus_Happy(t *testing.T) {
	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	getStatus := model.StatusTodo

	repo := RepoJsonFileMap{fileName: fileName}
	todoStatus, err := repo.GetByStatus(nil, getStatus)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if len(todos) != len(todoStatus) {
		t.Errorf("unexpected length expectedTodosMap != length todosMap")
		return
	}

	for i := range todos {
		expected := todos[i]
		actuals := todoStatus[i]
		if expected != actuals {
			t.Errorf("unexpected value, expecting='%+v', got='%+v'", expected, actuals)
			return
		}
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestGetStatus_Err(t *testing.T) {
	// status.IsValid()
}

func TestUpdateData_Happy(t *testing.T) {
	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	id := "1"
	newData := "oneone"

	repo := RepoJsonFileMap{fileName: fileName}
	_, err = repo.UpdateData(nil, id, newData)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	newTodoMap, err := readDecode(repo.fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if newData != newTodoMap[id].Data {
		t.Errorf("expected data: `%s` but got `%s`", newData, newTodoMap[id].Data)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestUpdateStatus_Happy(t *testing.T) {
	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	id := "1"
	newStatus := model.StatusDone

	repo := RepoJsonFileMap{fileName: fileName}
	_, err = repo.UpdateStatus(nil, id, newStatus)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	todosMapStatus, err := readDecode(repo.fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if newStatus != todosMapStatus[id].Status {
		t.Errorf("expected status: `%s` but got `%s`", newStatus, todosMapStatus[id].Data)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestRemove_Happy(t *testing.T) {
	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	id := "1"

	repo := RepoJsonFileMap{fileName: fileName}
	_, err = repo.Remove(nil, id)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	newTodosMap, err := readDecode(repo.fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	_, ok := newTodosMap[id]
	if ok {
		t.Errorf("unexpected found id: `%s` but found id: `%s`", id, newTodosMap[id].Id)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}

func TestRemove_IdErr(t *testing.T) {
	expectedErr := "no id"

	todosMap := make(map[string]model.Todo)
	todos := []model.Todo{
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
			Status: model.StatusTodo,
		},
	}

	for i := range todos {
		value := todos[i]
		key := value.Id
		todosMap[key] = value
	}

	err := writeEncode(fileName, todosMap)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	id := "66"

	repo := RepoJsonFileMap{fileName: fileName}
	_, err = repo.Remove(nil, id)
	if err == nil {
		t.Errorf("expected err but got nil")
		return
	}

	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected %s but got %s", expectedErr, err)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected writefile err: `%s`", err)
		return
	}

}
