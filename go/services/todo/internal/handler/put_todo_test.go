package handler_test

import (
	"context"
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

func Test_PutTodo(t *testing.T) {
	type fields struct {
		mockTodoUpdater *mock_usecase.MockTodoUpdater
	}

	type args struct {
		ctx context.Context
		req *todo_v1.PutTodoRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *todo_v1.PutTodoResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"Put Todo successfully": {
			prepare: func(f *fields) {
				f.mockTodoUpdater.
					EXPECT().
					Update(gomock.Any(), &input.TodoUpdater{
						ID:          todo.TodoID(1),
						Task:        cast.Ptr("updated task"),
						Description: cast.Ptr("updated description"),
						Status:      cast.Ptr(todo.InProcess),
					}).
					Return((*output.TodoUpdater)(&todo.Todo{
						ID:          todo.TodoID(1),
						UserID:      todo.UserID(1),
						Task:        "updated task",
						Description: cast.Ptr("updated description"),
						Status:      todo.InProcess,
					}), nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.PutTodoRequest{
					UserAttributes: &todo_v1.UserAttributes{
						UserId: 1,
					},
					TodoId:      1,
					Task:        "updated task",
					Description: "updated description",
					Status:      todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
				},
			},
			expected: &todo_v1.PutTodoResponse{
				Todo: &todo_common_v1.Todo{
					Id:          1,
					Task:        "updated task",
					Description: "updated description",
					Status:      todo_common_v1.TodoStatus_TODO_STATUS_INPROCESS,
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockTodoUpdater := mock_usecase.NewMockTodoUpdater(ctrl)

			tt.prepare(&fields{
				mockTodoUpdater: mockTodoUpdater,
			})

			svc, err := handler.NewTodoService(
				nil,
				validate,
				requestBinder,
				nil, nil, nil,
				mockTodoUpdater,
				nil, nil, nil,
			)
			if err != nil {
				t.Fatalf("NewTodoService() error = %v", err)
			}

			got, err := svc.PutTodo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got,
				protocmp.Transform(),
				protocmp.IgnoreFields(&todo_common_v1.Todo{}, "user_id", "created_at", "updated_at"),
			); diff != "" {
				t.Errorf("PutTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
