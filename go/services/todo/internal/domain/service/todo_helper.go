package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type todoHelper struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoHelper(
	todoQueriesGateway gateway.TodoQueriesGateway,
) TodoHelper {
	return &todoHelper{
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (h *todoHelper) CanAccessTodo(
	ctx context.Context,
	userID todo.UserID,
	todoID todo.TodoID,
) (bool, error) {
	todo, err := h.todoQueriesGateway.GetTodo(ctx, todoID, userID)
	if err != nil {
		return false, app_errors.NewInternalError(
			"todoHelper.CanAccessTodo",
			err,
		)
	}
	if todo == nil {
		return false, app_errors.NewNotFoundError(
			"todoHelper.CanAccessTodo", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}
	if todo.UserID != userID {
		return false, app_errors.NewAuthZError(
			"todoHelper.CanAccessTodo", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.AuthZJaMessage},
		)
	}
	return true, nil
}
