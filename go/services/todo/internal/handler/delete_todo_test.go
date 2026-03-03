package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	mock_usecase "github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/mock"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

func Test_DeleteTodo(t *testing.T) {
	type fields struct {
		mockTodoDeleter *mock_usecase.MockTodoDeleter
	}

	type args struct {
		ctx context.Context
		req *todo_v1.DeleteTodoRequest
	}

	type testcase struct {
		prepare  func(f *fields)
		args     args
		expected *todo_v1.DeleteTodoResponse
		wantErr  bool
	}

	validate := validator.New()
	requestBinder := requestbinder.NewRequestBinder(validate)

	testTables := map[string]testcase{
		"Delete Todo successfully": {
			prepare: func(f *fields) {
				f.mockTodoDeleter.
					EXPECT().
					SoftDelete(gomock.Any(), &input.TodoDeleter{
						ID: 1,
					}).
					Return(nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.DeleteTodoRequest{
					UserAttributes: &todo_v1.UserAttributes{
						UserId: 1,
					},
					TodoId: 1,
				},
			},
			expected: &todo_v1.DeleteTodoResponse{},
			wantErr:  false,
		},
		"Missing todo_id returns validation error": {
			prepare: func(f *fields) {
				f.mockTodoDeleter.EXPECT().SoftDelete(gomock.Any(), gomock.Any()).Times(0)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.DeleteTodoRequest{
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
				f.mockTodoDeleter.
					EXPECT().
					SoftDelete(gomock.Any(), &input.TodoDeleter{
						ID: 1,
					}).
					Return(errors.New("todo not found")).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				req: &todo_v1.DeleteTodoRequest{
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

			mockTodoDeleter := mock_usecase.NewMockTodoDeleter(ctrl)

			tt.prepare(&fields{
				mockTodoDeleter: mockTodoDeleter,
			})

			svc, err := handler.NewTodoService(
				nil,
				validate,
				requestBinder,
				nil, nil, nil, nil,
				mockTodoDeleter,
				nil, nil,
			)
			if err != nil {
				t.Fatalf("NewTodoService() error = %v", err)
			}

			got, err := svc.DeleteTodo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.expected, got, cmpopts.IgnoreUnexported(todo_v1.DeleteTodoResponse{})); diff != "" {
				t.Errorf("DeleteTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
