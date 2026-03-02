package output

import (
	"time"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type TodoGetter struct {
	ID          todo.TodoID
	Task        string
	Description *string
	Status      todo.TodoStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoLister struct {
	Todos []*todo.Todo
	Total int
}

type TodoCreator todo.Todo

type TodoUpdater todo.Todo

type UserGetter struct {
	ID       todo.UserID
	Username string
	Email    *string
}

type UserCreator struct {
	ID       todo.UserID
	Username string
	Email    *string
}

type UserRequestValidator todo.User
