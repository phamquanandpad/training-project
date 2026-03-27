package todo

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"

	todo_model "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
)

func Test_grpcTodoToModel(t *testing.T) {
	type args struct {
		pbTodo *todo_common_v1.Todo
	}

	type expected struct {
		todo *todo_model.Todo
	}

	now := time.Now()
	createdAt, updatedAt := timestamppb.New(now), timestamppb.New(now)

	testTables := map[string]struct {
		args     args
		expected expected
	}{
		"Convert gRPC Todo to model Todo successfully": {
			args: args{
				pbTodo: &todo_common_v1.Todo{
					Id:          1,
					UserId:      1,
					Task:        "todo task 1",
					Description: "todo description 1",
					Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
			},
			expected: expected{
				todo: &todo_model.Todo{
					ID:          todo_model.TodoID(1),
					UserID:      todo_model.UserID(1),
					Task:        "todo task 1",
					Description: cast.Ptr("todo description 1"),
					Status:      todo_model.Pending,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		"Convert nil gRPC Todo to nil model Todo": {
			args: args{
				pbTodo: nil,
			},
			expected: expected{
				todo: nil,
			},
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			out := grpcTodoToModel(tt.args.pbTodo)
			if diff := cmp.Diff(out, tt.expected.todo); diff != "" {
				t.Errorf("grpcTodoToModel() = got differs from expected: %s", diff)
			}
		})
	}
}

func Test_grpcTodoStatusToModelStatus(t *testing.T) {
	type args struct {
		statusGrpc todo_common_v1.TodoStatus
	}

	type expected struct {
		statusModel todo_model.TodoStatus
	}

	testTables := map[string]struct {
		args     args
		expected expected
	}{
		"Convert gRPC TodoStatus DONE to model TodoStatus Done": {
			args: args{
				statusGrpc: todo_common_v1.TodoStatus_TODO_STATUS_DONE,
			},
			expected: expected{
				statusModel: todo_model.Done,
			},
		},
		"Convert gRPC TodoStatus PENDING to model TodoStatus Pending": {
			args: args{
				statusGrpc: todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
			},
			expected: expected{
				statusModel: todo_model.Pending,
			},
		},
		"Convert gRPC TodoStatus INPROCESS to model TodoStatus InProcess": {
			args: args{
				statusGrpc: todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
			},
			expected: expected{
				statusModel: todo_model.InProcess,
			},
		},
		"Convert unknown gRPC TodoStatus to model TodoStatus Pending": {
			args: args{
				statusGrpc: 999,
			},
			expected: expected{
				statusModel: todo_model.Pending,
			},
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			got := grpcTodoStatusToModelStatus(tt.args.statusGrpc)
			if got != tt.expected.statusModel {
				t.Errorf("grpcTodoStatusToModelStatus() = got %v, want %v", got, tt.expected.statusModel)
			}
		})
	}
}

func Test_todoStatusToGrpcStatus(t *testing.T) {
	type args struct {
		statusModel todo_model.TodoStatus
	}

	type expected struct {
		statusGrpc todo_common_v1.TodoStatus
	}

	testTables := map[string]struct {
		args     args
		expected expected
	}{
		"Convert model TodoStatus Done to gRPC TodoStatus DONE": {
			args: args{
				statusModel: todo_model.Done,
			},
			expected: expected{
				statusGrpc: todo_common_v1.TodoStatus_TODO_STATUS_DONE,
			},
		},
		"Convert model TodoStatus Pending to gRPC TodoStatus PENDING": {
			args: args{
				statusModel: todo_model.Pending,
			},
			expected: expected{
				statusGrpc: todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
			},
		},
		"Convert model TodoStatus InProcess to gRPC TodoStatus INPROCESS": {
			args: args{
				statusModel: todo_model.InProcess,
			},
			expected: expected{
				statusGrpc: todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
			},
		},
		"Convert unknown model TodoStatus to gRPC TodoStatus PENDING": {
			args: args{
				statusModel: 999,
			},
			expected: expected{
				statusGrpc: todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
			},
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			got := todoStatusToGrpcStatus(tt.args.statusModel)
			if got != tt.expected.statusGrpc {
				t.Errorf("todoStatusToGrpcStatus() = got %v, want %v", got, tt.expected.statusGrpc)
			}
		})
	}
}

func Test_grpcListTodosToModels(t *testing.T) {
	type args struct {
		todosGrpc []*todo_common_v1.Todo
	}

	type expected struct {
		todosModel []*todo_model.Todo
	}

	now := time.Now()
	createdAt, updatedAt := timestamppb.New(now), timestamppb.New(now)

	testTables := map[string]struct {
		args     args
		expected expected
	}{
		"Convert list of gRPC Todos to list of model Todos successfully": {
			args: args{
				todosGrpc: []*todo_common_v1.Todo{
					{
						Id:          1,
						UserId:      1,
						Task:        "todo task 1",
						Description: "todo description 1",
						Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
					{
						Id:          2,
						UserId:      1,
						Task:        "todo task 2",
						Description: "todo description 2",
						Status:      todo_common_v1.TodoStatus_TODO_STATUS_DONE,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
				},
			},
			expected: expected{
				todosModel: []*todo_model.Todo{
					{
						ID:          todo_model.TodoID(1),
						UserID:      todo_model.UserID(1),
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo_model.Pending,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          todo_model.TodoID(2),
						UserID:      todo_model.UserID(1),
						Task:        "todo task 2",
						Description: cast.Ptr("todo description 2"),
						Status:      todo_model.Done,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
		},
		"Convert empty list of gRPC Todos to empty list of model Todos": {
			args: args{
				todosGrpc: []*todo_common_v1.Todo{},
			},
			expected: expected{
				todosModel: []*todo_model.Todo{},
			},
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			got := grpcListTodosToModels(tt.args.todosGrpc)
			if diff := cmp.Diff(tt.expected.todosModel, got); diff != "" {
				t.Errorf("grpcListTodosToModels() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
