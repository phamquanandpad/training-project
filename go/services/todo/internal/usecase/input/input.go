package input

import (
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type UserAttributes struct {
	UserID todo.UserID `json:"user_id" validate:"required"`
}

type TodoGetter struct {
	ID todo.TodoID `json:"todo_id" validate:"required"`
}

type TodoLister struct {
	UserAttributes UserAttributes `json:"user_attributes"`
	Offset         int64          `json:"offset"`
	Limit          int64          `json:"limit"`
}

type TodoCreator struct {
	Task        string
	Description *string
	Status      todo.TodoStatus
}

type TodoUpdater struct {
	ID          todo.TodoID `json:"todo_id" validate:"required"`
	Task        *string
	Description *string
	Status      *todo.TodoStatus
}

type TodoDeleter struct {
	ID todo.TodoID `json:"todo_id" validate:"required"`
}

type UserGetter struct {
	UserID todo.UserID `json:"user_id" validate:"required"`
}

type UserCreator struct {
	User todo.User `json:"user"`
}

type UserRequestValidator struct {
	UserID todo.UserID
}
