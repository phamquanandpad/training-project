package output

import (
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type TodoGetter todo.Todo

type TodoLister struct {
	Todos []*todo.Todo
	Total int
}

type TodoCreator todo.Todo

type TodoUpdater todo.Todo

type UserGetter todo.User

type UserCreator todo.User

type UserRequestValidator todo.User
