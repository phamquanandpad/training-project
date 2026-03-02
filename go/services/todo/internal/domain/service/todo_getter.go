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

type todoGetter struct {
	dbConnBinder       gateway.Binder
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoGetter(
	dbConnBinder gateway.Binder,
	todoQueriesGateway gateway.TodoQueriesGateway,
) usecase.TodoGetter {
	return &todoGetter{
		dbConnBinder:       dbConnBinder,
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (s todoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
	ctx = s.dbConnBinder.Bind(ctx)
	currentUser := todo.ExtractUser(ctx)

	todo, err := s.todoQueriesGateway.GetTodo(ctx, in.ID, currentUser.ID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoGetter.GetTodo",
			err,
		)
	}

	if todo == nil {
		return nil, app_errors.NewNotFoundError(
			"service.GetTodo", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}

	if todo.UserID != currentUser.ID {
		return nil, app_errors.NewAuthZError(
			"todoGetter.GetTodo", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.AuthZJaMessage},
		)
	}

	return &output.TodoGetter{
		ID:          todo.ID,
		Task:        todo.Task,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}
