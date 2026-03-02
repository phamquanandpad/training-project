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

func (h *todoService) ListTodos(
	ctx context.Context,
	req *todo_v1.ListTodosRequest,
) (*todo_v1.ListTodosResponse, error) {
	var in input.TodoLister

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	todos, err := h.todoLister.List(ctx, &in)
	if err != nil {
		return nil, err
	}

	return toListTodosResponse(todos), nil
}

// func toListTodosInput(in *todo_v1.ListTodosRequest) *input.TodoLister {
// 	return &input.TodoLister{
// 		UserAttributes: input.UserAttributes{
// 			UserID: todo.UserID(in.GetUserAttributes().GetUserId()),
// 		},
// 		Offset: int(in.GetOffset()),
// 		Limit:  int(in.GetLimit()),
// 	}
// }

func toListTodosResponse(out *output.TodoLister) *todo_v1.ListTodosResponse {
	todos := make([]*todo_common_v1.Todo, len(out.Todos))
	for i, t := range out.Todos {
		todos[i] = &todo_common_v1.Todo{
			Id:          int64(t.ID),
			Task:        t.Task,
			Description: cast.Value(t.Description),
			Status:      todo_common_v1.TodoStatus(t.Status),
			CreatedAt:   timestamppb.New(t.CreatedAt),
			UpdatedAt:   timestamppb.New(t.UpdatedAt),
		}
	}

	return &todo_v1.ListTodosResponse{
		Todos: todos,
		Total: int64(out.Total),
	}
}
