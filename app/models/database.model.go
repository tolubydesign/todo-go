package models

type ToDoItem struct {
	Id   string `json:"id,omitempty" validate:"uuid"`
	Task string `json:"task" validate:"required"`
}
