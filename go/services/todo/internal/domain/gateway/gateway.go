package gateway

import (
	"context"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type Binder interface {
	Bind(context.Context) context.Context
}

type TodoQueriesGateway interface {
	GetTodo(ctx context.Context, todoID todo.TodoID) (*todo.Todo, error)
	ListTodos(ctx context.Context) ([]*todo.Todo, error)
}

type TodoCommandsGateway interface {
	CreateTodo(ctx context.Context, newTodo todo.NewTodo) (*todo.Todo, error)
	UpdateTodo(ctx context.Context, todoID todo.TodoID, updateTodo todo.UpdateTodo) (*todo.Todo, error)
	SoftDeleteTodo(ctx context.Context, todoID todo.TodoID) error
}

type UserQueriesGateway interface {
	GetUser(ctx context.Context, userID int64) (*todo.User, error)
}

type UserCommandsGateway interface {
	CreateUser(ctx context.Context, newUser todo.NewUser) (*todo.User, error)
}
