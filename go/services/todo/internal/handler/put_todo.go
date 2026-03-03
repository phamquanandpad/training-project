package handler

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/mapper"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

func (h *todoService) PutTodo(
	ctx context.Context,
	req *todo_v1.PutTodoRequest,
) (*todo_v1.PutTodoResponse, error) {
	var in input.TodoUpdater

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	out, err := h.todoUpdater.Update(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &todo_v1.PutTodoResponse{
		Todo: mapper.ToTodoGRPCResponse((*todo.Todo)(out)),
	}, nil
}
