package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
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

	return toPostTodoResponse(out), nil
}

func toPostTodoInput(in *todo_v1.PostTodoRequest) *input.TodoCreator {
	return &input.TodoCreator{
		Task:        in.GetTask(),
		Description: cast.Ptr(in.GetDescription()),
		Status:      todo.TodoStatus(in.GetStatus()),
	}
}

func toPostTodoResponse(out *output.TodoCreator) *todo_v1.PostTodoResponse {
	return &todo_v1.PostTodoResponse{
		Todo: &todo_common_v1.Todo{
			Id:          int64(out.ID),
			Task:        out.Task,
			Description: cast.Value(out.Description),
			Status:      todo_common_v1.TodoStatus(out.Status),
			CreatedAt:   timestamppb.New(out.CreatedAt),
			UpdatedAt:   timestamppb.New(out.UpdatedAt),
		},
	}
}
