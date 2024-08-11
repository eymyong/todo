package model

import "errors"

type Todo struct {
	Id     string `json:"id"`
	Data   string `json:"data"`
	Status Status `json:"status"`
}

type Status string

func (s Status) MarshalBinary() (data []byte, err error) {
	if s == "" {
		return nil, errors.New("status is empty string")
	}

	return []byte(s), nil
}

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
	Id     string `json:"id"`
	Data   int    `json:"data"`
	Status Status `json:"status"`
}
