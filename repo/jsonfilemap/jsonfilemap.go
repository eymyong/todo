package jsonfilemap

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eymyong/todo/model"
	"github.com/eymyong/todo/repo"
)

type RepoJsonFileMap struct {
	fileName string
}

func readDecode(fileName string) (map[string]model.Todo, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read jsonfile: %w", err)
	}

	if len(b) == 0 {
		return make(map[string]model.Todo), nil
	}

	todos := make(map[string]model.Todo)
	err = json.Unmarshal(b, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func writeEncode(fileName string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = os.WriteFile(fileName, b, 0664)
	if err != nil {
		return fmt.Errorf("failed to writefile jsonfile: %w", err)
	}

	return nil
}
func (j *RepoJsonFileMap) Add(todo model.Todo) error {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return err
	}

	todoMap[todo.Id] = todo
	writeEncode(j.fileName, todoMap)

	return nil
}

func (j *RepoJsonFileMap) GetAll() ([]model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return []model.Todo{}, err
	}

	todoList := make([]model.Todo, len(todoMap))

	i := 0
	for _, todo := range todoMap {
		todoList[i] = todo
		i++
	}

	return todoList, nil
}

func (j *RepoJsonFileMap) Get(id string) (model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	// if todoMap = map[]  จะให้ return err ออกเลยและแจ้งว่า `no data`
	fmt.Println(todoMap)

	todo, ok := todoMap[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("no id: %s", id)
	}

	return todo, nil
}

func (j *RepoJsonFileMap) GetStatus(status model.Status) ([]model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return []model.Todo{}, err
	}

	// if todoMap = map[]  จะให้ return err ออกเลยและแจ้งว่า `no data`

	checkStatus := status.IsValid()
	if !checkStatus {
		return []model.Todo{}, fmt.Errorf("bad status: `%s`", status)
	}

	newTodos := []model.Todo{}

	for _, todo := range todoMap {
		if todo.Status == status {
			newTodos = append(newTodos, todo)
		}
	}

	return newTodos, nil
}

func (j *RepoJsonFileMap) UpdateData(id string, newData string) (model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	// if todoMap = map[]  จะให้ return err ออกเลยและแจ้งว่า `no data`

	old, ok := todoMap[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("no id: %s", id)
	}

	copy := old
	copy.Data = newData
	todoMap[id] = copy

	err = writeEncode(j.fileName, todoMap)
	if err != nil {
		return model.Todo{}, err
	}

	return old, nil
}

func (j *RepoJsonFileMap) UpdateStatus(id string, newStatus model.Status) (model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	// if todoMap = map[]  จะให้ return err ออกเลยและแจ้งว่า `no data`

	old, ok := todoMap[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("no id: %s", id)
	}

	copy := old
	copy.Status = newStatus
	todoMap[id] = copy

	err = writeEncode(j.fileName, todoMap)
	if err != nil {
		return model.Todo{}, err
	}

	return old, nil
}

func (j *RepoJsonFileMap) Remove(id string) (model.Todo, error) {
	todoMap, err := readDecode(j.fileName)
	if err != nil {
		return model.Todo{}, err
	}

	// if todoMap = map[]  จะให้ return err ออกเลยและแจ้งว่า `no data`

	todo, ok := todoMap[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("no id: %s", id)
	}

	//delete(todoMap, id)
	// ถ้าใช้ delete ไม่ต้องทำข้างล่าง

	todos := []model.Todo{}
	for _, v := range todoMap {
		if v.Id == id {
			continue
		}

		todos = append(todos, v)
	}

	newTodoMap := make(map[string]model.Todo)
	for i := range todos {
		value := todos[i]
		key := value.Id
		newTodoMap[key] = value
	}

	err = writeEncode(j.fileName, newTodoMap)
	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func New(fileName string) repo.Repository {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil || len(fileBytes) == 0 {
		err := os.WriteFile(fileName, []byte("{}"), os.ModePerm)
		if err != nil {
			panic("failed to write empty array to init file: " + err.Error())
		}
	}

	return &RepoJsonFileMap{
		fileName: fileName,
	}
}
