package jsonfile

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/eymyong/todo/model"
	"github.com/google/uuid"
)

const fileName = "mock/test_foo.json"
const fileNameErr = "error.json"

func TestReadDecodeErrorFileName(t *testing.T) {
	expectedErr := "failed to read jsonfile"
	_, err := readDecode(fileNameErr)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
	}
}

// TODO
func TestRdadDecodeUnmarshalErr(t *testing.T) {
	todos, err := readDecode(fileName)
	if len(todos) != 0 {
		t.Errorf("expected error but got nil")
	}

	expectedErr := "failed to unmarshal"

	expectedTodos := []model.TestTodo{
		{
			Id:     1,
			Data:   "kuy",
			Status: model.StatusTodo,
		},
	}

	byteTodos, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err)
		return
	}

	err = os.WriteFile(fileName, byteTodos, 0664)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err)
		return
	}

	_, err = readDecode(fileName)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err)
		return
	}

}

func TestReadDecodeHappy(t *testing.T) {
	expectedTodos := []model.Todo{
		{
			Id:   "2",
			Data: "kuy",
		},
	}

	b, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	if todos[0].Id != expectedTodos[0].Id {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodos[0].Id, todos[0].Id)
	}

	if todos[0].Data != expectedTodos[0].Data {
		t.Errorf("expected data: '%s' but got '%s'", expectedTodos[0].Id, todos[0].Id)
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

}

func TestAddHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := model.Todo{
		Id:     uuid.NewString(),
		Data:   "yong",
		Status: model.StatusTodo,
	}

	err := repo.Add(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}
	if todos[0].Id != expectedTodos.Id {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodos.Id, todos[0].Id)
	}
	if todos[0].Data != expectedTodos.Data {
		t.Errorf("expected data: '%s' but got '%s'", expectedTodos.Data, todos[0].Data)
	}

	if todos[0].Status != expectedTodos.Status {
		t.Errorf("expected status: '%s' but got '%s'", expectedTodos.Status, todos[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}
}

func TestAddFileError(t *testing.T) {
	repo := RepoJsonFile{
		fileName: "aaa",
	}
	expectedAddErr := "failed to add jsonfile"
	expectedTodos := model.Todo{
		Id:     uuid.NewString(),
		Data:   "yong",
		Status: model.StatusTodo,
	}

	err := repo.Add(expectedTodos)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedAddErr) {
		t.Errorf("expected '%s' but got '%s'", expectedAddErr, err.Error())
	}
}

func TestUpdateDataHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := []model.Todo{{
		Id:     "1",
		Data:   "yong",
		Status: model.StatusTodo,
	}}

	b, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	newData := "pak"

	_, err = repo.UpdateData(expectedTodos[0].Id, newData)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	if todos[0].Id != expectedTodos[0].Id {
		t.Errorf("expected id: `%s` but got `%s", expectedTodos[0].Id, todos[0].Id)
	}

	if newData != todos[0].Data {
		t.Errorf("expected data: '%s' but got '%s'", newData, todos[0].Data)
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

// TODO
func TestUpdateDataErr(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}
	expectedAddErr := "failed to updatedata jsonfile"
	expectedTodos := model.Todo{
		Id:     "1",
		Data:   "yong",
		Status: model.StatusTodo,
	}
	newData := "pak"

	_, err := repo.UpdateData(expectedTodos.Id, newData)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedAddErr) {
		t.Errorf("expected '%s' but got '%s'", expectedAddErr, err.Error())
	}

}

// TODO
func TestUpdateStatusHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := []model.Todo{{
		Id:     "1",
		Data:   "yong",
		Status: model.StatusTodo,
	}}

	newStatus := model.StatusDone

	b, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	_, err = repo.UpdateStatus(expectedTodos[0].Id, newStatus)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	if newStatus != todos[0].Status {
		t.Errorf("expected data: '%s' but got '%s'", newStatus, todos[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

}

func TestUpdateStatusErr(t *testing.T) {

}

func TestDeleteHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := []model.Todo{
		{
			Id:     "1",
			Data:   "kuy",
			Status: "TODO",
		},
		{
			Id:     "2",
			Data:   "hee",
			Status: "DONE",
		},
		{
			Id:     "3",
			Data:   "hum",
			Status: "DONE",
		},
	}

	b, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	_, err = repo.Remove(expectedTodos[1].Id)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	if len(todos) != len(expectedTodos)-1 {
		t.Errorf("unexpected length todos != '%v'", len(expectedTodos)-1)
		return
	}

	for _, v := range todos {
		if v.Id == expectedTodos[1].Id {
			t.Errorf("unexpected id: '%s'", v.Id)
		}
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

func TestDeleteErr(t *testing.T) {

}

func TestGetAll(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}
	expectedTodos := []model.Todo{
		{
			Id:   "2",
			Data: "kuy",
		},
	}

	b, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	todos, err := repo.GetAll()
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	if expectedTodos[0].Id != todos[0].Id {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodos[0].Id, todos[0].Id)
	}

	if expectedTodos[0].Data != todos[0].Data {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodos[0].Data, todos[0].Data)
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}
}

func TestGetByID(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := []model.Todo{
		{
			Id:   "2",
			Data: "kuy",
		},
	}

	b, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	todo, err := repo.Get(expectedTodos[0].Id)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	if todo.Data != expectedTodos[0].Data {
		t.Errorf("expected data: '%s' but got '%s'", expectedTodos[0].Data, todo.Data)
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

}

func TesyGetSatatus(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := []model.Todo{
		{
			Id:     "1",
			Data:   "kuy",
			Status: "TODO",
		},
		{
			Id:     "2",
			Data:   "hee",
			Status: "DONE",
		},
		{
			Id:     "3",
			Data:   "hum",
			Status: "DONE",
		},
	}

	allStatus := []string{"TODO", "DONE"}

	b, err := json.Marshal(expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	for _, status := range allStatus {
		if status == "DONE" {
			todos, err := repo.GetStatus("DONE")
			if err != nil {
				t.Errorf("unexpected err: %s", err.Error())
			}

			for _, todo := range todos {
				if todo.Status != "DONE" {
					t.Errorf("expected status: '%s' but got status: '%s'", "DONE", todo.Status)
				}
				return
			}

			err = os.WriteFile(fileName, []byte{}, os.ModePerm)
			if err != nil {
				t.Errorf("unexpected err: %s", err.Error())
			}

		}

		todos, err := repo.GetStatus("TODO")
		if err != nil {
			t.Errorf("unexpected err: %s", err.Error())
		}

		for _, todo := range todos {
			if todo.Status != "TODO" {
				t.Errorf("expected status: '%s' but got status: '%s'", "TODO", todo.Status)
			}
			return
		}

		err = os.WriteFile(fileName, []byte{}, os.ModePerm)
		if err != nil {
			t.Errorf("unexpected err: %s", err.Error())
		}

	}
}
