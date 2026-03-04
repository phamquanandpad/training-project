package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

type todoDeleter struct {
	dbConnBinder        gateway.Binder
	todoCommandsGateway gateway.TodoCommandsGateway
	todoHelper          TodoHelper
}

func NewTodoDeleter(
	dbConnBinder gateway.Binder,
	todoCommandsGateway gateway.TodoCommandsGateway,
	todoHelper TodoHelper,
) usecase.TodoDeleter {
	return &todoDeleter{
		dbConnBinder:        dbConnBinder,
		todoCommandsGateway: todoCommandsGateway,
		todoHelper:          todoHelper,
	}
}

func (s todoDeleter) SoftDelete(ctx context.Context, in *input.TodoDeleter) error {
	ctx = s.dbConnBinder.Bind(ctx)
	currentUser := todo.ExtractUser(ctx)

	canAccess, err := s.todoHelper.CanAccessTodo(ctx, currentUser.ID, in.ID)
	if !canAccess || err != nil {
		return err
	}

	err = s.todoCommandsGateway.SoftDeleteTodo(ctx, in.ID, currentUser.ID)
	if err != nil {
		return app_errors.NewInternalError(
			"todoDeleter.Delete",
			err,
		)
	}

	return nil
}
