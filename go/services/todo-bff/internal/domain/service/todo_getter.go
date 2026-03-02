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

type todoGetter struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoGetter(
	todoQueriesGateway gateway.TodoQueriesGateway,
) usecase.TodoGetter {
	return &todoGetter{
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (s todoGetter) GetTodo(
	ctx context.Context,
	in *input.TodoGetter,
) (*output.TodoGetter, error) {
	todo, err := s.todoQueriesGateway.GetTodo(
		ctx,
		todo.UserAttributes{
			UserID: in.UserID,
		},
		in.ID,
	)
	if todo == nil {
		return nil, app_errors.NewNotFoundError(
			"todoGetter.GetTodo",
			fmt.Errorf("todo not found"),
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}
	if err != nil {
		return nil, app_errors.NewInternalError(
			"todoGetter.GetTodo",
			err,
		)
	}

	return &output.TodoGetter{
		Todo: todo,
	}, nil
}
