package service

import (
	"context"
	"fmt"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type todoUpdater struct {
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoUpdater(
	todoCommandsGateway gateway.TodoCommandsGateway,
) usecase.TodoUpdater {
	return &todoUpdater{
		todoCommandsGateway: todoCommandsGateway,
	}
}

func (s todoUpdater) UpdateTodo(
	ctx context.Context,
	in *input.TodoUpdater,
) (*output.TodoUpdater, error) {
	updatedTodo, err := s.todoCommandsGateway.UpdateTodo(
		ctx,
		todo.UserAttributes{
			UserID: in.UserID,
		},
		in.ID,
		todo.UpdateTodo{
			Task:        in.Task,
			Description: in.Description,
			Status:      in.Status,
		},
	)
	if updatedTodo == nil {
		return nil, app_errors.NewNotFoundError(
			"todoUpdater.UpdateTodo",
			fmt.Errorf("todo not found"),
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoUpdater.UpdateTodo",
			err,
		)
	}

	return &output.TodoUpdater{
		Todo: updatedTodo,
	}, nil
}
