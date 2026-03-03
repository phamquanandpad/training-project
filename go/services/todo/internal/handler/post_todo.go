package handler

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/mapper"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

func (h *todoService) PostTodo(
	ctx context.Context,
	req *todo_v1.PostTodoRequest,
) (*todo_v1.PostTodoResponse, error) {
	var in input.TodoCreator

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	out, err := h.todoCreator.Create(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &todo_v1.PostTodoResponse{
		Todo: mapper.ToTodoGRPCResponse((*todo.Todo)(out)),
	}, nil
}
