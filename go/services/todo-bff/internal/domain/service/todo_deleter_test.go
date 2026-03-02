package service_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway/mock"

	todo_model "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
)

type PrepareTodoDeleterFields struct {
	ctx                     context.Context
	mockTodoCommandsGateway *mock_gateway.MockTodoCommandsGateway
}

type TodoDeleterArgs struct {
	ctx context.Context
	in  *input.TodoDeleter
}

type TodoDeleterTestcase struct {
	prepare func(f *PrepareTodoDeleterFields)
	args    TodoDeleterArgs
	wantErr bool
}

func Test_todoDeleter_DeleteTodo(t *testing.T) {
	t.Parallel()

	testTables := map[string]TodoDeleterTestcase{
		"Delete Todo successfully": {
			prepare: func(f *PrepareTodoDeleterFields) {
				f.mockTodoCommandsGateway.
					EXPECT().
					DeleteTodo(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.TodoID(1)).
					Return(nil).
					Times(1)
			},
			args: TodoDeleterArgs{
				ctx: context.Background(),
				in: &input.TodoDeleter{
					UserID: todo_model.UserID(1),
					ID:     todo_model.TodoID(1),
				},
			},
			wantErr: false,
		},
		"Fail to delete Todo when gateway returns error": {
			prepare: func(f *PrepareTodoDeleterFields) {
				f.mockTodoCommandsGateway.
					EXPECT().
					DeleteTodo(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.TodoID(1)).
					Return(errors.New("gateway error")).
					Times(1)
			},
			args: TodoDeleterArgs{
				ctx: context.Background(),
				in: &input.TodoDeleter{
					UserID: todo_model.UserID(1),
					ID:     todo_model.TodoID(1),
				},
			},
			wantErr: true,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockTodoCommandsGateway := mock_gateway.NewMockTodoCommandsGateway(ctrl)

			f := &PrepareTodoDeleterFields{
				ctx:                     context.Background(),
				mockTodoCommandsGateway: mockTodoCommandsGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewTodoDeleter(mockTodoCommandsGateway)

			err := svc.DeleteTodo(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("todoDeleter.DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("todoDeleter.DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
