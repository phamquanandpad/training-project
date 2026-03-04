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

type todoLister struct {
	dbConnBinder       gateway.Binder
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoLister(
	dbConnBinder gateway.Binder,
	todoQueriesGateway gateway.TodoQueriesGateway,
) usecase.TodoLister {
	return &todoLister{
		dbConnBinder:       dbConnBinder,
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (s todoLister) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	ctx = s.dbConnBinder.Bind(ctx)
	currentUser := todo.ExtractUser(ctx)

	todos, total, err := s.todoQueriesGateway.ListTodos(
		ctx,
		currentUser.ID,
		int(in.Limit),
		int(in.Offset),
	)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoLister.List",
			err,
		)
	}

	return &output.TodoLister{
		Todos: todos,
		Total: total,
	}, nil
}
