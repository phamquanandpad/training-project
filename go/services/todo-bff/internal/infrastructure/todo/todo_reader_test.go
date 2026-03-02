package todo

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	mock_todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1/mock"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	todo_model "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
)

func Test_todoReader_GetTodo(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockTodoServiceClient *mock_todo_v1.MockTodoServiceClient
	}

	type args struct {
		ctx            context.Context
		userAttributes todo_model.UserAttributes
		todoID         todo_model.TodoID
	}

	testTables := map[string]struct {
		prepare  func(f *fields)
		args     args
		expected *todo_model.Todo
		wantErr  bool
	}{
		"Get Todo successfully": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					GetTodo(gomock.Any(), &todo_v1.GetTodoRequest{
						UserAttributes: &todo_v1.UserAttributes{
							UserId: 1,
						},
						TodoId: 1,
					}).
					Return(&todo_v1.GetTodoResponse{
						Todo: &todo_common_v1.Todo{
							Id:          1,
							UserId:      1,
							Task:        "todo task 1",
							Description: "todo description 1",
							Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
						},
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				userAttributes: todo_model.UserAttributes{
					UserID: todo_model.UserID(1),
				},
				todoID: todo_model.TodoID(1),
			},
			expected: &todo_model.Todo{
				ID:          todo_model.TodoID(1),
				UserID:      todo_model.UserID(1),
				Task:        "todo task 1",
				Description: cast.Ptr("todo description 1"),
				Status:      todo_model.Pending,
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			f := &fields{
				mockTodoServiceClient: mock_todo_v1.NewMockTodoServiceClient(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			reader := NewTodoReader(f.mockTodoServiceClient)
			actual, err := reader.GetTodo(tt.args.ctx, tt.args.userAttributes, tt.args.todoID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("TodoReader.GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("TodoReader.GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(todo_model.Todo{}, "CreatedAt", "UpdatedAt"),
			}

			if diff := cmp.Diff(tt.expected, actual, ignoreFieldsOpts...); diff != "" {
				t.Errorf("TodoReader.GetTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_todoReader_ListTodos(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockTodoServiceClient *mock_todo_v1.MockTodoServiceClient
	}

	type args struct {
		ctx            context.Context
		userAttributes todo_model.UserAttributes
		limit          int
		offset         int
	}

	type expected struct {
		todos []*todo_model.Todo
		total int
	}

	limit := 10
	offset := 0

	now := time.Now()
	createdAt := timestamppb.New(now)
	updatedAt := timestamppb.New(now)

	testTables := map[string]struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
	}{
		"List Todos successfully": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					ListTodos(gomock.Any(), &todo_v1.ListTodosRequest{
						UserAttributes: &todo_v1.UserAttributes{
							UserId: 1,
						},
						Limit:  cast.Ptr(int64(limit)),
						Offset: cast.Ptr(int64(offset)),
					}).
					Return(&todo_v1.ListTodosResponse{
						Todos: []*todo_common_v1.Todo{
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
								Status:      todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
								CreatedAt:   createdAt,
								UpdatedAt:   updatedAt,
							},
						},
						Total: 2,
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				userAttributes: todo_model.UserAttributes{
					UserID: todo_model.UserID(1),
				},
				limit:  limit,
				offset: offset,
			},
			expected: expected{
				todos: []*todo_model.Todo{
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
						Status:      todo_model.InProcess,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
				total: 2,
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			f := &fields{
				mockTodoServiceClient: mock_todo_v1.NewMockTodoServiceClient(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			reader := NewTodoReader(f.mockTodoServiceClient)
			actual, total, err := reader.ListTodos(tt.args.ctx, tt.args.userAttributes, tt.args.limit, tt.args.offset)
			if tt.wantErr {
				if err == nil {
					t.Errorf("TodoReader.ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("TodoReader.ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected.todos, actual); diff != "" {
				t.Errorf("TodoReader.ListTodos() mismatch (-want +got):\n%s", diff)
			}
			if tt.expected.total != total {
				t.Errorf("TodoReader.ListTodos() total = %v, want %v", total, tt.expected.total)
			}
		})
	}
}
