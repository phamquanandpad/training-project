package todo

import (
	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
)

func grpcTodoToModel(todo_grpc *todo_common_v1.Todo) *todo.Todo {
	if todo_grpc == nil {
		return nil
	}

	return &todo.Todo{
		ID:          todo.TodoID(todo_grpc.Id),
		UserID:      todo.UserID(todo_grpc.UserId),
		Task:        todo_grpc.Task,
		Description: cast.Ptr(todo_grpc.Description),
		Status:      grpcTodoStatusToModelStatus(todo_grpc.Status),
		CreatedAt:   todo_grpc.CreatedAt.AsTime(),
		UpdatedAt:   todo_grpc.UpdatedAt.AsTime(),
	}
}

func grpcListTodosToModels(todos_grpc []*todo_common_v1.Todo) []*todo.Todo {
	todos := make([]*todo.Todo, 0, len(todos_grpc))
	for _, todo_grpc := range todos_grpc {
		todos = append(todos, grpcTodoToModel(todo_grpc))
	}
	return todos
}

func grpcTodoStatusToModelStatus(status_grpc todo_common_v1.TodoStatus) todo.TodoStatus {
	switch status_grpc {
	case todo_common_v1.TodoStatus_TODO_STATUS_DONE:
		return todo.Done
	case todo_common_v1.TodoStatus_TODO_STATUS_PENDING:
		return todo.Pending
	case todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS:
		return todo.InProcess
	default:
		return todo.Pending
	}
}

func todoStatusToGrpcStatus(status todo.TodoStatus) todo_common_v1.TodoStatus {
	switch status {
	case todo.Done:
		return todo_common_v1.TodoStatus_TODO_STATUS_DONE
	case todo.Pending:
		return todo_common_v1.TodoStatus_TODO_STATUS_PENDING
	case todo.InProcess:
		return todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS
	default:
		return todo_common_v1.TodoStatus_TODO_STATUS_PENDING
	}
}
