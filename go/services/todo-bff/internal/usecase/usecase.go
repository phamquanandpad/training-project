package usecase

import (
	"context"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type UserLogin interface {
	Login(
		ctx context.Context,
		in *input.UserLogin,
	) (*output.UserLogin, error)
}

type UserRegister interface {
	Register(
		ctx context.Context,
		in *input.UserRegister,
	) error
}

type TokenRefresh interface {
	RefreshToken(
		ctx context.Context,
		in *input.TokenRefresh,
	) (*output.TokenRefresh, error)
}

type TokenVerify interface {
	VerifyToken(
		ctx context.Context,
		in *input.TokenVerify,
	) (*output.TokenVerify, error)
}

type TodoGetter interface {
	GetTodo(
		ctx context.Context,
		in *input.TodoGetter,
	) (*output.TodoGetter, error)
}

type TodoLister interface {
	ListTodos(
		ctx context.Context,
		in *input.TodoLister,
	) (*output.TodoLister, error)
}

type TodoCreator interface {
	CreateTodo(
		ctx context.Context,
		in *input.TodoCreator,
	) (*output.TodoCreator, error)
}

type TodoUpdater interface {
	UpdateTodo(
		ctx context.Context,
		in *input.TodoUpdater,
	) (*output.TodoUpdater, error)
}

type TodoDeleter interface {
	DeleteTodo(
		ctx context.Context,
		in *input.TodoDeleter,
	) error
}
