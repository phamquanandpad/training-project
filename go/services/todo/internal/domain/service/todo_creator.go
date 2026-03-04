package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type todoCreator struct {
	dbConnBinder        gateway.Binder
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoCreator(
	dbConnBinder gateway.Binder,
	todoCommandsGateway gateway.TodoCommandsGateway,
) usecase.TodoCreator {
	return &todoCreator{
		dbConnBinder:        dbConnBinder,
		todoCommandsGateway: todoCommandsGateway,
	}
}

func (s todoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	ctx = s.dbConnBinder.Bind(ctx)
	currentUser := todo.ExtractUser(ctx)

	newTodo := todo.NewTodo{
		UserID:      currentUser.ID,
		Task:        in.Task,
		Description: in.Description,
		Status:      in.Status,
	}

	todo, err := s.todoCommandsGateway.CreateTodo(ctx, newTodo)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoCreator.Create",
			err,
		)
	}

	return &output.TodoCreator{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Task:        todo.Task,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}
