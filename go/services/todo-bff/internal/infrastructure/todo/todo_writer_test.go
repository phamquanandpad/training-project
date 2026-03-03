package todo

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"

	mock_todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1/mock"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	todo_model "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
)

func Test_todoWriter_Create(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockTodoServiceClient *mock_todo_v1.MockTodoServiceClient
	}

	type args struct {
		ctx            context.Context
		userAttributes todo_model.UserAttributes
		newTodo        todo_model.NewTodo
	}

	testTables := map[string]struct {
		prepare  func(f *fields)
		args     args
		expected *todo_model.Todo
		wantErr  bool
	}{
		"Create Todo successfully": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					PostTodo(gomock.Any(), &todo_v1.PostTodoRequest{
						UserAttributes: &todo_v1.UserAttributes{
							UserId: 1,
						},
						Task:        "todo task 1",
						Description: "todo description 1",
						Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
					}).
					Return(&todo_v1.PostTodoResponse{
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
				newTodo: todo_model.NewTodo{
					Task:        "todo task 1",
					Description: cast.Ptr("todo description 1"),
					Status:      todo_model.Pending,
				},
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
		"Fail when gRPC client returns error on CreateTodo": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					PostTodo(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("rpc error: internal server error")).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				userAttributes: todo_model.UserAttributes{
					UserID: todo_model.UserID(1),
				},
				newTodo: todo_model.NewTodo{
					Task:   "todo task 1",
					Status: todo_model.Pending,
				},
			},
			expected: nil,
			wantErr:  true,
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

			writer := NewTodoWriter(f.mockTodoServiceClient)
			actual, err := writer.Create(tt.args.ctx, tt.args.userAttributes, tt.args.newTodo)
			if tt.wantErr {
				if err == nil {
					t.Errorf("TodoWriter.CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("TodoWriter.CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(todo_model.Todo{}, "ID", "CreatedAt", "UpdatedAt"),
			}

			if diff := cmp.Diff(tt.expected, actual, ignoreFieldsOpts...); diff != "" {
				t.Errorf("TodoWriter.CreateTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_todoWriter_Update(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockTodoServiceClient *mock_todo_v1.MockTodoServiceClient
	}

	type args struct {
		ctx            context.Context
		userAttributes todo_model.UserAttributes
		todoID         todo_model.TodoID
		updateTodo     todo_model.UpdateTodo
	}

	testTables := map[string]struct {
		prepare  func(f *fields)
		args     args
		expected *todo_model.Todo
		wantErr  bool
	}{
		"Update Todo successfully": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					PutTodo(gomock.Any(), &todo_v1.PutTodoRequest{
						UserAttributes: &todo_v1.UserAttributes{
							UserId: 1,
						},
						TodoId:      1,
						Task:        "todo task updated",
						Description: "todo description updated",
						Status:      todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
					}).
					Return(&todo_v1.PutTodoResponse{
						Todo: &todo_common_v1.Todo{
							Id:          1,
							UserId:      1,
							Task:        "todo task updated",
							Description: "todo description updated",
							Status:      todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
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
				updateTodo: todo_model.UpdateTodo{
					Task:        cast.Ptr("todo task updated"),
					Description: cast.Ptr("todo description updated"),
					Status:      cast.Ptr(todo_model.InProcess),
				},
			},
			expected: &todo_model.Todo{
				ID:          todo_model.TodoID(1),
				UserID:      todo_model.UserID(1),
				Task:        "todo task updated",
				Description: cast.Ptr("todo description updated"),
				Status:      todo_model.InProcess,
			},
			wantErr: false,
		},
		"Fail when gRPC client returns error on UpdateTodo": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					PutTodo(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("rpc error: todo not found")).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				userAttributes: todo_model.UserAttributes{
					UserID: todo_model.UserID(1),
				},
				todoID: todo_model.TodoID(999),
				updateTodo: todo_model.UpdateTodo{
					Task: cast.Ptr("updated task"),
				},
			},
			expected: nil,
			wantErr:  true,
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

			writer := NewTodoWriter(f.mockTodoServiceClient)
			actual, err := writer.Update(tt.args.ctx, tt.args.userAttributes, tt.args.todoID, tt.args.updateTodo)
			if tt.wantErr {
				if err == nil {
					t.Errorf("TodoWriter.UpdateTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("TodoWriter.UpdateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(todo_model.Todo{}, "CreatedAt", "UpdatedAt"),
			}

			if diff := cmp.Diff(tt.expected, actual, ignoreFieldsOpts...); diff != "" {
				t.Errorf("TodoWriter.UpdateTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_todoWriter_Delete(t *testing.T) {
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
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		"Delete Todo successfully": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					DeleteTodo(gomock.Any(), &todo_v1.DeleteTodoRequest{
						UserAttributes: &todo_v1.UserAttributes{
							UserId: 1,
						},
						TodoId: 1,
					}).
					Return(&todo_v1.DeleteTodoResponse{}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				userAttributes: todo_model.UserAttributes{
					UserID: todo_model.UserID(1),
				},
				todoID: todo_model.TodoID(1),
			},
			wantErr: false,
		},
		"Fail when gRPC client returns error on DeleteTodo": {
			prepare: func(f *fields) {
				f.mockTodoServiceClient.
					EXPECT().
					DeleteTodo(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("rpc error: todo not found")).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				userAttributes: todo_model.UserAttributes{
					UserID: todo_model.UserID(1),
				},
				todoID: todo_model.TodoID(999),
			},
			wantErr: true,
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

			writer := NewTodoWriter(f.mockTodoServiceClient)
			err := writer.Delete(tt.args.ctx, tt.args.userAttributes, tt.args.todoID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("TodoWriter.DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("TodoWriter.DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
