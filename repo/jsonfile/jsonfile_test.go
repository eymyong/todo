package jsonfile

import (
	"encoding/json"
	"github.com/eymyong/todo/model"
	"os"
	"strings"
	"testing"
)

const fileName = "mock/test_foo.json"

func TestReadDecodeErrorFileName(t *testing.T) {
	expectedErr := "failed to read jsonfile"
	_, err := readDecode("error.json")
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
	}
}

func TestReadDecodeErrorUnmarshal(t *testing.T) {
	expectedErr := "unexpected end of JSON input"
	_, err := readDecode("mock/test_foo_unmarshal_err.json")
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
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
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
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
	}

	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	todos, err := repo.GetAll()
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	if expectedTodos[0].Id == todos[0].Id {
		//TODO:
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

}
