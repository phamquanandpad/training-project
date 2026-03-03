package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	mock_usecase "github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

func Test_ListTodos(t *testing.T) {
	type fields struct {
		mockTodoLister *mock_usecase.MockTodoLister
	}

	type args struct {
		ctx context.Context
		req *todo_v1.ListTodosRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *todo_v1.ListTodosResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"List Todos successfully": {
			prepare: func(f *fields) {
				f.mockTodoLister.
					EXPECT().
					List(gomock.Any(), &input.TodoLister{
						UserAttributes: input.UserAttributes{UserID: 1},
						Offset:         0,
						Limit:          10,
					}).
					Return(&output.TodoLister{
						Todos: []*todo.Todo{
							{
								ID:          todo.TodoID(1),
								UserID:      todo.UserID(1),
								Task:        "todo task 1",
								Description: cast.Ptr("todo description 1"),
								Status:      todo.Pending,
							},
						},
						Total: 1,
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.ListTodosRequest{
					UserAttributes: &todo_v1.UserAttributes{
						UserId: 1,
					},
					Limit: cast.Ptr(int64(10)),
				},
			},
			expected: &todo_v1.ListTodosResponse{
				Todos: []*todo_common_v1.Todo{
					{
						Id:          1,
						UserId:      1,
						Task:        "todo task 1",
						Description: "todo description 1",
						Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
					},
				},
				Total: 1,
			},
			wantErr: false,
		},
		"Usecase returns error": {
			prepare: func(f *fields) {
				f.mockTodoLister.
					EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.ListTodosRequest{
					UserAttributes: &todo_v1.UserAttributes{
						UserId: 1,
					},
					Limit: cast.Ptr(int64(10)),
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

			mockTodoLister := mock_usecase.NewMockTodoLister(ctrl)

			tt.prepare(&fields{
				mockTodoLister: mockTodoLister,
			})

			svc, err := handler.NewTodoService(
				nil,
				validate,
				requestBinder,
				nil,
				mockTodoLister,
				nil, nil, nil, nil, nil,
			)
			if err != nil {
				t.Fatalf("NewTodoService() error = %v", err)
			}

			got, err := svc.ListTodos(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got,
				protocmp.Transform(),
				protocmp.IgnoreFields(&todo_common_v1.Todo{}, "created_at", "updated_at"),
			); diff != "" {
				t.Errorf("ListTodos() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
