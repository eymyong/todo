package model

type Todo struct {
	Id     string `json:"id"`
	Data   string `json:"data"`
	Status Status `json:"status"`
}

type Status string

const (
	StatusTodo Status = "TODO"
	StatusDone Status = "DONE"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusDone, "":
		return true
	}

	return false
}

type TestTodo struct {
	Id     int    `json:"id"`
	Data   string `json:"data"`
	Status Status `json:"status"`
}
