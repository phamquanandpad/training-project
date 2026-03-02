package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
)

type todoDeleter struct {
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoDeleter(
	todoCommandsGateway gateway.TodoCommandsGateway,
) usecase.TodoDeleter {
	return &todoDeleter{
		todoCommandsGateway: todoCommandsGateway,
	}
}

func (s todoDeleter) DeleteTodo(
	ctx context.Context,
	in *input.TodoDeleter,
) error {
	err := s.todoCommandsGateway.DeleteTodo(
		ctx,
		todo.UserAttributes{
			UserID: in.UserID,
		},
		in.ID,
	)
	if err != nil {
		return app_errors.NewInternalError(
			"todoDeleter.DeleteTodo",
			err,
		)
	}
	return err
}
