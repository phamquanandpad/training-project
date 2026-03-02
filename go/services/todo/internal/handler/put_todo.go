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

	return toPutTodoResponse(out), nil
}

// func toPutTodoInput(in *todo_v1.PutTodoRequest) *input.TodoUpdater {
// 	return &input.TodoUpdater{
// 		ID:          todo.TodoID(in.GetTodoId()),
// 		Task:        cast.Ptr(in.GetTask()),
// 		Description: cast.Ptr(in.GetDescription()),
// 		Status:      cast.Ptr(todo.TodoStatus(in.GetStatus())),
// 	}
// }

func toPutTodoResponse(out *output.TodoUpdater) *todo_v1.PutTodoResponse {
	return &todo_v1.PutTodoResponse{
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
