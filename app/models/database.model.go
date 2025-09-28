package models

type ToDoItem struct {
	Id       string `json:"id,omitempty" validate:"uuid"`
	Task     string `json:"task" validate:"required"`
	Complete bool   `json:"complete" validate:"required"`
}
