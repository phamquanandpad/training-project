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

type todoLister struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoLister(
	todoQueriesGateway gateway.TodoQueriesGateway,
) usecase.TodoLister {
	return &todoLister{
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (s todoLister) ListTodos(
	ctx context.Context,
	in *input.TodoLister,
) (*output.TodoLister, error) {
	todos, total, err := s.todoQueriesGateway.ListTodos(
		ctx,
		todo.UserAttributes{
			UserID: in.UserID,
		},
		in.Limit,
		in.Offset,
	)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoLister.ListTodos",
			err,
		)
	}

	return &output.TodoLister{
		Todos:      todos,
		TotalCount: total,
	}, nil
}
