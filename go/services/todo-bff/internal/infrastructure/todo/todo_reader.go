package todo

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
)

type todoReader struct {
	todoServiceClient todo_v1.TodoServiceClient
}

func NewTodoReader(client todo_v1.TodoServiceClient) gateway.TodoQueriesGateway {
	return &todoReader{
		todoServiceClient: client,
	}
}

func (r *todoReader) Get(
	ctx context.Context,
	userAttributes todo.UserAttributes,
	todoID todo.TodoID,
) (*todo.Todo, error) {
	req := &todo_v1.GetTodoRequest{
		UserAttributes: &todo_v1.UserAttributes{
			UserId: int64(userAttributes.UserID),
		},
		TodoId: int64(todoID),
	}

	res, err := r.todoServiceClient.GetTodo(ctx, req)
	if err != nil {
		return nil, app_errors.GrpcStatusToAppError(err)
	}

	return grpcTodoToModel(res.Todo), nil
}

func (r *todoReader) List(
	ctx context.Context,
	userAttributes todo.UserAttributes,
	limit int,
	offset int,
) ([]*todo.Todo, int, error) {
	req := &todo_v1.ListTodosRequest{
		UserAttributes: &todo_v1.UserAttributes{
			UserId: int64(userAttributes.UserID),
		},
		Limit:  cast.Ptr(int64(limit)),
		Offset: cast.Ptr(int64(offset)),
	}

	res, err := r.todoServiceClient.ListTodos(ctx, req)
	if err != nil {
		return nil, 0, app_errors.GrpcStatusToAppError(err)
	}

	todos := grpcListTodosToModels(res.Todos)

	return todos, int(res.Total), nil
}
