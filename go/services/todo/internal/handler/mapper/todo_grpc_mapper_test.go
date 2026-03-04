package mapper_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/mapper"
)

func Test_ToTodoGRPCResponse(t *testing.T) {
	type args struct {
		todo *todo.Todo
	}

	now := time.Now()

	testTables := map[string]struct {
		args     args
		expected *todo_common_v1.Todo
	}{
		"Convert Todo to gRPC response successfully": {
			args: args{
				todo: &todo.Todo{
					ID:          todo.TodoID(1),
					UserID:      todo.UserID(2),
					Task:        "todo task 1",
					Description: cast.Ptr("todo description 1"),
					Status:      todo.Pending,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			expected: &todo_common_v1.Todo{
				Id:          1,
				UserId:      2,
				Task:        "todo task 1",
				Description: "todo description 1",
				Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
				CreatedAt:   timestamppb.New(now),
				UpdatedAt:   timestamppb.New(now),
			},
		},
		"Return nil when Todo is nil": {
			args: args{
				todo: nil,
			},
			expected: nil,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			out := mapper.ToTodoGRPCResponse(tt.args.todo)
			if diff := cmp.Diff(out, tt.expected, protocmp.Transform()); diff != "" {
				t.Errorf("ToTodoGRPCResponse() differs from expected: %s", diff)
			}
		})
	}
}

func Test_ToTodosGRPCResponse(t *testing.T) {
	type args struct {
		todos []*todo.Todo
	}

	now := time.Now()

	testTables := map[string]struct {
		args     args
		expected []*todo_common_v1.Todo
	}{
		"Convert list of Todos to gRPC response successfully": {
			args: args{
				todos: []*todo.Todo{
					{
						ID:          todo.TodoID(1),
						UserID:      todo.UserID(1),
						Task:        "todo task 1",
						Description: cast.Ptr("description 1"),
						Status:      todo.Pending,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          todo.TodoID(2),
						UserID:      todo.UserID(1),
						Task:        "todo task 2",
						Description: cast.Ptr("description 2"),
						Status:      todo.InProcess,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			expected: []*todo_common_v1.Todo{
				{
					Id:          1,
					UserId:      1,
					Task:        "todo task 1",
					Description: "description 1",
					Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
					CreatedAt:   timestamppb.New(now),
					UpdatedAt:   timestamppb.New(now),
				},
				{
					Id:          2,
					UserId:      1,
					Task:        "todo task 2",
					Description: "description 2",
					Status:      todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
					CreatedAt:   timestamppb.New(now),
					UpdatedAt:   timestamppb.New(now),
				},
			},
		},
		"Convert empty list of Todos to empty gRPC response": {
			args: args{
				todos: []*todo.Todo{},
			},
			expected: []*todo_common_v1.Todo{},
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			out := mapper.ToTodosGRPCResponse(tt.args.todos)
			if diff := cmp.Diff(out, tt.expected, protocmp.Transform()); diff != "" {
				t.Errorf("ToTodosGRPCResponse() differs from expected: %s", diff)
			}
		})
	}
}

func Test_ToTodoStatusGRPCResponse(t *testing.T) {
	type args struct {
		status todo.TodoStatus
	}

	testTables := map[string]struct {
		args     args
		expected todo_common_v1.TodoStatus
	}{
		"Convert Pending status to TODO_STATUS_PENDING": {
			args:     args{status: todo.Pending},
			expected: todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
		},
		"Convert InProcess status to TODO_STATUS_INPROCESS": {
			args:     args{status: todo.InProcess},
			expected: todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
		},
		"Convert Done status to TODO_STATUS_DONE": {
			args:     args{status: todo.Done},
			expected: todo_common_v1.TodoStatus_TODO_STATUS_DONE,
		},
		"Convert unknown status to TODO_STATUS_PENDING (default)": {
			args:     args{status: todo.TodoStatus(99)},
			expected: todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			out := mapper.ToTodoStatusGRPCResponse(tt.args.status)
			if diff := cmp.Diff(out, tt.expected); diff != "" {
				t.Errorf("ToTodoStatusGRPCResponse() differs from expected: %s", diff)
			}
		})
	}
}
