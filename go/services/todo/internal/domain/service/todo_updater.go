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

type todoUpdater struct {
	dbConnBinder        gateway.Binder
	todoCommandsGateway gateway.TodoCommandsGateway
	todoHelper          TodoHelper
}

func NewTodoUpdater(
	dbConnBinder gateway.Binder,
	todoCommandsGateway gateway.TodoCommandsGateway,
	todoHelper TodoHelper,
) usecase.TodoUpdater {
	return &todoUpdater{
		dbConnBinder:        dbConnBinder,
		todoCommandsGateway: todoCommandsGateway,
		todoHelper:          todoHelper,
	}
}

func (s todoUpdater) Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error) {
	ctx = s.dbConnBinder.Bind(ctx)
	currentUser := todo.ExtractUser(ctx)

	canAccess, err := s.todoHelper.CanAccessTodo(ctx, currentUser.ID, in.ID)
	if !canAccess || err != nil {
		return nil, err
	}

	updateTodo := todo.UpdateTodo{
		Task:        in.Task,
		Description: in.Description,
		Status:      in.Status,
	}

	todo, err := s.todoCommandsGateway.UpdateTodo(ctx, in.ID, currentUser.ID, updateTodo)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoUpdater.Update",
			err,
		)
	}

	return &output.TodoUpdater{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Task:        todo.Task,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}
