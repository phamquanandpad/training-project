package todo

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
)

type todoWriter struct {
	todoServiceClient todo_v1.TodoServiceClient
}

func NewTodoWriter(client todo_v1.TodoServiceClient) gateway.TodoCommandsGateway {
	return &todoWriter{
		todoServiceClient: client,
	}
}

func (w *todoWriter) Create(
	ctx context.Context,
	userAttributes todo.UserAttributes,
	newTodo todo.NewTodo,
) (*todo.Todo, error) {
	req := &todo_v1.PostTodoRequest{
		UserAttributes: &todo_v1.UserAttributes{
			UserId: int64(userAttributes.UserID),
		},
		Task:        newTodo.Task,
		Description: cast.Value(newTodo.Description),
		Status:      todoStatusToGrpcStatus(newTodo.Status),
	}

	res, err := w.todoServiceClient.PostTodo(ctx, req)
	if err != nil {
		return nil, app_errors.GrpcStatusToAppError(err)
	}

	return grpcTodoToModel(res.Todo), nil
}

func (w *todoWriter) Update(
	ctx context.Context,
	userAttributes todo.UserAttributes,
	todoID todo.TodoID,
	updateTodo todo.UpdateTodo,
) (*todo.Todo, error) {
	req := &todo_v1.PutTodoRequest{
		UserAttributes: &todo_v1.UserAttributes{
			UserId: int64(userAttributes.UserID),
		},
		TodoId:      int64(todoID),
		Task:        cast.Value(updateTodo.Task),
		Description: cast.Value(updateTodo.Description),
		Status:      todoStatusToGrpcStatus(cast.Value(updateTodo.Status)),
	}

	res, err := w.todoServiceClient.PutTodo(ctx, req)
	if err != nil {
		return nil, app_errors.GrpcStatusToAppError(err)
	}

	return grpcTodoToModel(res.Todo), nil
}

func (w *todoWriter) Delete(
	ctx context.Context,
	userAttributes todo.UserAttributes,
	todoID todo.TodoID,
) error {
	req := &todo_v1.DeleteTodoRequest{
		UserAttributes: &todo_v1.UserAttributes{
			UserId: int64(userAttributes.UserID),
		},
		TodoId: int64(todoID),
	}

	_, err := w.todoServiceClient.DeleteTodo(ctx, req)
	if err != nil {
		return app_errors.GrpcStatusToAppError(err)
	}

	return nil
}
