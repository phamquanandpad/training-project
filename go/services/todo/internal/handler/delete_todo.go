package handler

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

func (h *todoService) DeleteTodo(
	ctx context.Context,
	req *todo_v1.DeleteTodoRequest,
) (*todo_v1.DeleteTodoResponse, error) {
	var in input.TodoDeleter

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	err := h.todoDeleter.SoftDelete(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &todo_v1.DeleteTodoResponse{}, nil
}
