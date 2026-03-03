package handler

import (
	"context"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/mapper"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

func (h *todoService) ListTodos(
	ctx context.Context,
	req *todo_v1.ListTodosRequest,
) (*todo_v1.ListTodosResponse, error) {
	var in input.TodoLister

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	out, err := h.todoLister.List(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &todo_v1.ListTodosResponse{
		Todos: mapper.ToTodosGRPCResponse(out.Todos),
		Total: int64(out.Total),
	}, nil
}
