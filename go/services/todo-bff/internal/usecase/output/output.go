package output

import (
	auth_models "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/auth"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
)

type UserLogin struct {
	Tokens *auth_models.Tokens
}

type UserRegister struct{}

type TokenRefresh struct {
	AccessToken *auth_models.AccessToken
}

type TokenVerify struct {
	UserID *auth_models.UserID
}

type TodoGetter struct {
	Todo *todo.Todo
}

type TodoLister struct {
	Todos      []*todo.Todo
	TotalCount int
}

type TodoCreator struct {
	Todo *todo.Todo
}

type TodoUpdater struct {
	Todo *todo.Todo
}

type TodoDeleter struct{}
