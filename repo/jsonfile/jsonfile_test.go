package jsonfile

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/eymyong/todo/model"
)

/*
go test
go test -v
*/
const fileName = "mock/test_foo.json"
const fileNameErr = "error.json"

func TestReadDecodeHappy(t *testing.T) {
	todosNoData, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	if len(todosNoData) != 0 {
		t.Errorf("expected length `todosNoData` = 0, but got length `todosNoData` != 0")
		return
	}

	expectedTodos := makeTodos()

	err = writeEncode(fileName, expectedTodos)
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

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

// ถ้ามี file "error.json" อยู่แล้ว จะError ต้องมาทำให้ถูก
func TestReadDecodeFileName_Err(t *testing.T) {
	expectedErr := "failed to read jsonfile"
	_, err := readDecode("expected err")
	if err != nil {
		t.Errorf("expected error but got nil:")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
	}
}

func TestReadDecodeUnmarshal_Err(t *testing.T) {
	expectedErr := "failed to unmarshal"
	todos, err := readDecode(fileNameErr)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
	}
	//ไม่มีข้อมูลไฟล์ readDecode return []model.Todo{} >> readDecode does not perform Unmarshal >> จึงต้องมาเช็คแยก
	if len(todos) == 0 {
		expectedErrUnmarshal := "unexpected end of JSON input"
		b, err := os.ReadFile(fileNameErr)
		if err != nil {
			t.Errorf("unexpected err: `%s`", err)
		}

		var u model.Todo
		err = json.Unmarshal(b, &u)
		if err == nil {
			t.Errorf("expected err but got nil")
		}

		// error contain expected error ?
		if !strings.Contains(err.Error(), expectedErrUnmarshal) {
			t.Errorf("expected '%s' but got '%s'", expectedErrUnmarshal, err.Error())
			return
		}

		return
	}
	// make new field in struct
	expectedTodos := makeTodosTest()

	err = writeEncode(fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err.Error())
		return
	}

	_, err = readDecode(fileName)
	if err == nil {
		t.Errorf("expected error but got nil")
		return
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

func TestWriteEncodeHappy(t *testing.T) {
	expectedTodos := makeTodos()

	err := writeEncode(fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	if expectedTodos[0].Id != todos[0].Id {
		t.Errorf("expected id: '%s' but got '%s'", expectedTodos[0].Id, todos[0].Id)
	}

	if expectedTodos[0].Data != todos[0].Data {
		t.Errorf("expected data: '%s' but got '%s'", expectedTodos[0].Data, todos[0].Data)
	}

	if expectedTodos[0].Status != todos[0].Status {
		t.Errorf("expected status: '%s' but got '%s'", expectedTodos[0].Status, todos[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

// TODO ต้องใส่ข้อมูลให้ผิด
// func TestWriteEncodeWriteFileErr(t *testing.T) {
// 	expectedTodos := makeTodosTest()
// 	expectedErr := "failed to marshal"

// 	err := writeEncode(fileNameErr, expectedTodos)
// 	if err == nil {
// 		t.Errorf("expected error but got nil")
// 	}

// 	// error contain expected error ?
// 	if !strings.Contains(err.Error(), expectedErr) {
// 		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
// 	}

// 	err = os.WriteFile(fileNameErr, []byte{}, 0664)
// 	if err != nil {
// 		t.Errorf("unexpected err: %s", err.Error())
// 	}

// }

func TestAddHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := makeTodos()
	newTodo := expectedTodos[0]

	err := repo.Add(newTodo)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	todos, err := readDecode(repo.fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}
	if todos[0].Id != newTodo.Id {
		t.Errorf("expected id: '%s' but got '%s'", newTodo.Id, todos[0].Id)
	}
	if todos[0].Data != newTodo.Data {
		t.Errorf("expected data: '%s' but got '%s'", newTodo.Data, todos[0].Data)
	}

	if todos[0].Status != newTodo.Status {
		t.Errorf("expected status: '%s' but got '%s'", newTodo.Status, todos[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}
}

func TestAddFileError(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileNameErr,
	}
	expectedErr := "failed to add"

	expectedTodos := makeTodos()

	dataToAdd := expectedTodos[0]

	err := repo.Add(dataToAdd)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// error contain expected error ?
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected '%s' but got '%s'", expectedErr, err.Error())
	}
}

func TestUpdateDataHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := makeTodos()
	updateTo := expectedTodos[0]

	err := writeEncode(repo.fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err)
		return
	}

	newData := "pak"

	_, err = repo.UpdateData(updateTo.Id, newData)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if updateTo.Id != todos[0].Id {
		t.Errorf("expected id: `%s` but got `%s", updateTo.Id, todos[0].Id)
	}

	if newData != todos[0].Data {
		t.Errorf("expected data: '%s' but got '%s'", newData, todos[0].Data)
	}

	if updateTo.Status != todos[0].Status {
		t.Errorf("expected status: '%s' but got '%s'", updateTo.Status, todos[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

func TestUpdateStatusHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := makeTodos()

	updateTo := expectedTodos[0]

	newStatus := model.StatusDone

	err := writeEncode(repo.fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err.Error())
		return
	}

	_, err = repo.UpdateStatus(updateTo.Id, newStatus)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

	if newStatus != todos[0].Status {
		t.Errorf("expected status: '%s' but got '%s'", newStatus, todos[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}

}

func TestDeleteHappy(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := makeTodos()

	deleteToID := expectedTodos[1].Id

	lengthExpectedTodos := len(expectedTodos)

	err := writeEncode(repo.fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err.Error())
		return
	}

	_, err = repo.Remove(deleteToID)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	todos, err := readDecode(fileName)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	lengthTodos := len(todos)

	if lengthTodos != lengthExpectedTodos-1 {
		//t.Errorf("expected length-todos = '%v' but got length-todos != '%v'", lengthTodos, lengthTodos)
		t.Errorf("unexpected length-todos != '%v'", lengthTodos)
		return
	}

	for _, v := range todos {
		if deleteToID == v.Id {
			t.Errorf("unexpected found id: '%s'", deleteToID)
		}
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

func TestGetAll(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}
	expectedTodos := makeTodos()

	err := writeEncode(repo.fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err.Error())
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
		t.Errorf("expected data: '%s' but got '%s'", expectedTodos[0].Data, todos[0].Data)
	}

	if expectedTodos[0].Status != todos[0].Status {
		t.Errorf("expected status: '%s' but got '%s'", expectedTodos[0].Status, todos[0].Status)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}
}

func TestGet(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := makeTodos()
	get := expectedTodos[0]

	err := writeEncode(repo.fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpect err: `%s`", err.Error())
		return
	}

	todo, err := repo.Get(get.Id)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	if get.Data != todo.Data {
		t.Errorf("expected data: '%s' but got '%s'", expectedTodos[0].Data, todo.Data)
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

func TestGetSatatus(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := makeTodos()

	const TODO = model.StatusTodo
	const DONE = model.StatusDone

	err := writeEncode(repo.fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err.Error())
		return
	}

	var allStatus model.Status
	switch allStatus {
	case TODO:
		todosStatusTodo, err := repo.GetStatus(TODO)
		if err != nil {
			t.Errorf("unexpected err: `%s`", err.Error())
		}

		for _, v := range todosStatusTodo {
			if TODO != v.Status {
				t.Errorf("expected status: `%s` but got status: `%s`", TODO, v.Status)
			}
		}
	case "":
		fallthrough

	default:
		todosStatusDone, err := repo.GetStatus(DONE)
		if err != nil {
			t.Errorf("unexpected err: `%s`", err.Error())
		}

		for _, v := range todosStatusDone {
			if DONE != v.Status {
				t.Errorf("expected status: `%s` but got status: `%s`", DONE, v.Status)
			}
		}

	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}

func TestGetSatatus2(t *testing.T) {
	repo := RepoJsonFile{
		fileName: fileName,
	}

	expectedTodos := makeTodos()

	statusTODO := model.StatusTodo
	statusDONE := model.StatusDone

	err := writeEncode(repo.fileName, expectedTodos)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err.Error())
		return
	}

	todosStatusTodo, err := repo.GetStatus(statusTODO)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err.Error())
	}

	for _, v := range todosStatusTodo {
		if statusTODO != v.Status {
			t.Errorf("expected status: `%s` but got status: `%s`", statusTODO, v.Status)
		}
	}

	todosStatusDone, err := repo.GetStatus(statusDONE)
	if err != nil {
		t.Errorf("unexpected err: `%s`", err.Error())
	}

	for _, v := range todosStatusDone {
		if statusDONE != v.Status {
			t.Errorf("expected status: `%s` but got status: `%s`", statusDONE, v.Status)
		}
	}

	err = os.WriteFile(fileName, []byte{}, 0664)
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

}
