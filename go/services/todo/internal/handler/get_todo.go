package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

func (h *todoService) GetTodo(
	ctx context.Context,
	req *todo_v1.GetTodoRequest,
) (*todo_v1.GetTodoResponse, error) {
	var in input.TodoGetter

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	todo, err := h.todoGetter.Get(ctx, &in)
	if err != nil {
		return nil, err
	}

	return toGetTodoResponse(todo), nil
}

// func toGetTodoInput(in *todo_v1.GetTodoRequest) *input.TodoGetter {
// 	return &input.TodoGetter{
// 		ID: todo.TodoID(in.GetTodoId()),
// 	}
// }

func toGetTodoResponse(out *output.TodoGetter) *todo_v1.GetTodoResponse {
	return &todo_v1.GetTodoResponse{
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
