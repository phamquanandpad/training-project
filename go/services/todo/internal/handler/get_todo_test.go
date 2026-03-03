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

func Test_GetTodo(t *testing.T) {
	type fields struct {
		mockTodoGetter *mock_usecase.MockTodoGetter
	}

	type args struct {
		ctx context.Context
		req *todo_v1.GetTodoRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *todo_v1.GetTodoResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"Get Todo successfully": {
			prepare: func(f *fields) {
				f.mockTodoGetter.
					EXPECT().
					Get(gomock.Any(), &input.TodoGetter{
						ID: todo.TodoID(1),
					}).
					Return(&output.TodoGetter{
						ID:          todo.TodoID(1),
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo.Pending,
					}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.GetTodoRequest{
					UserAttributes: &todo_v1.UserAttributes{
						UserId: 1,
					},
					TodoId: 1,
				},
			},
			expected: &todo_v1.GetTodoResponse{
				Todo: &todo_common_v1.Todo{
					Id:          1,
					Task:        "todo task 1",
					Description: "todo description 1",
					Status:      todo_common_v1.TodoStatus_TODO_STATUS_PENDING,
				},
			},
			wantErr: false,
		},
		"Missing todo_id returns validation error": {
			prepare: func(f *fields) {
				f.mockTodoGetter.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.GetTodoRequest{
					UserAttributes: &todo_v1.UserAttributes{
						UserId: 1,
					},
					TodoId: 0,
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Usecase returns error": {
			prepare: func(f *fields) {
				f.mockTodoGetter.
					EXPECT().
					Get(gomock.Any(), &input.TodoGetter{
						ID: todo.TodoID(1),
					}).
					Return(nil, errors.New("todo not found")).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.GetTodoRequest{
					UserAttributes: &todo_v1.UserAttributes{
						UserId: 1,
					},
					TodoId: 1,
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

			mockTodoGetter := mock_usecase.NewMockTodoGetter(ctrl)

			tt.prepare(&fields{
				mockTodoGetter: mockTodoGetter,
			})

			svc, err := handler.NewTodoService(
				nil,
				validate,
				requestBinder,
				mockTodoGetter,
				nil, nil, nil, nil, nil, nil,
			)
			if err != nil {
				t.Fatalf("NewTodoService() error = %v", err)
			}

			got, err := svc.GetTodo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got,
				protocmp.Transform(),
				protocmp.IgnoreFields(&todo_common_v1.Todo{}, "created_at", "updated_at"),
			); diff != "" {
				t.Errorf("GetTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
