package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

func ToTodoGRPCResponse(todo *todo.Todo) *todo_common_v1.Todo {
	if todo == nil {
		return nil
	}

	return &todo_common_v1.Todo{
		Id:          todo.ID.Int64(),
		UserId:      todo.UserID.Int64(),
		Task:        todo.Task,
		Description: cast.Value(todo.Description),
		Status:      ToTodoStatusGRPCResponse(todo.Status),
		CreatedAt:   timestamppb.New(todo.CreatedAt),
		UpdatedAt:   timestamppb.New(todo.UpdatedAt),
	}
}

func ToTodosGRPCResponse(todos []*todo.Todo) []*todo_common_v1.Todo {
	result := make([]*todo_common_v1.Todo, len(todos))
	for i, todo := range todos {
		result[i] = ToTodoGRPCResponse(todo)
	}
	return result
}

func ToTodoStatusGRPCResponse(status todo.TodoStatus) todo_common_v1.TodoStatus {
	switch status {
	case todo.Pending:
		return todo_common_v1.TodoStatus_TODO_STATUS_PENDING
	case todo.InProcess:
		return todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS
	case todo.Done:
		return todo_common_v1.TodoStatus_TODO_STATUS_DONE
	default:
		return todo_common_v1.TodoStatus_TODO_STATUS_PENDING
	}
}
