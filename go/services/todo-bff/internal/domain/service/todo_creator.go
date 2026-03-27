package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type todoCreator struct {
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoCreator(
	todoCommandsGateway gateway.TodoCommandsGateway,
) usecase.TodoCreator {
	return &todoCreator{
		todoCommandsGateway: todoCommandsGateway,
	}
}

func (s todoCreator) CreateTodo(
	ctx context.Context,
	in *input.TodoCreator,
) (*output.TodoCreator, error) {
	todo, err := s.todoCommandsGateway.Create(
		ctx,
		todo.UserAttributes{
			UserID: in.UserID,
		},
		todo.NewTodo{
			Task:        in.Task,
			Description: in.Description,
			Status:      in.Status,
		},
	)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoCreator.CreateTodo",
			err,
		)
	}

	return &output.TodoCreator{
		Todo: todo,
	}, nil
}
