package input

import (
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
)

type UserLogin struct {
	Email    string
	Password string
}

type UserRegister struct {
	Username string
	Email    string
	Password string
}

type TokenRefresh struct {
	RefreshToken string
}

type TokenVerify struct {
	AccessToken string
}

type TodoGetter struct {
	UserID todo.UserID
	ID     todo.TodoID
}

type TodoLister struct {
	UserID todo.UserID
	Limit  int
	Offset int
}

type TodoCreator struct {
	UserID      todo.UserID
	Task        string
	Description *string
	Status      todo.TodoStatus
}

type TodoUpdater struct {
	UserID      todo.UserID
	ID          todo.TodoID
	Task        *string
	Description *string
	Status      *todo.TodoStatus
}

type TodoDeleter struct {
	UserID todo.UserID
	ID     todo.TodoID
}
