package service_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"
	mock_service "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service/mock"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
)

type PrepareDeleteTodoFields struct {
	ctx                     context.Context
	mockBinder              *mock_gateway.MockBinder
	mockTodoCommandsGateway *mock_gateway.MockTodoCommandsGateway
	mockTodoHelper          *mock_service.MockTodoHelper
}

type DeleteTodoArgs struct {
	ctx context.Context
	in  *input.TodoDeleter
}

type DeleteTodoTestcase struct {
	prepare func(f *PrepareDeleteTodoFields)
	args    DeleteTodoArgs
	wantErr bool
}

func Test_todoDeleter_SoftDelete(t *testing.T) {
	t.Parallel()

	existedUser := &todo.User{
		ID:       todo.UserID(1),
		Username: "user1",
	}

	testTables := map[string]DeleteTodoTestcase{
		"Soft delete Todo successfully": {
			prepare: func(f *PrepareDeleteTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoHelper.
					EXPECT().
					CanAccessTodo(f.ctx, todo.UserID(1), todo.TodoID(1)).
					Return(true, nil).
					Times(1)

				f.mockTodoCommandsGateway.
					EXPECT().
					SoftDeleteTodo(f.ctx, todo.TodoID(1), todo.UserID(1)).
					Return(nil).
					Times(1)
			},
			args: DeleteTodoArgs{
				ctx: context.Background(),
				in: &input.TodoDeleter{
					ID: todo.TodoID(1),
				},
			},
			wantErr: false,
		},
		"Internal error from CanAccessTodo": {
			prepare: func(f *PrepareDeleteTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoHelper.
					EXPECT().
					CanAccessTodo(f.ctx, todo.UserID(1), todo.TodoID(1)).
					Return(false, errors.New("database error")).
					Times(1)
			},
			args: DeleteTodoArgs{
				ctx: context.Background(),
				in: &input.TodoDeleter{
					ID: todo.TodoID(1),
				},
			},
			wantErr: true,
		},
		"Cannot access todo (unauthorized)": {
			prepare: func(f *PrepareDeleteTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoHelper.
					EXPECT().
					CanAccessTodo(f.ctx, todo.UserID(1), todo.TodoID(2)).
					Return(false, errors.New("unauthorized")).
					Times(1)
			},
			args: DeleteTodoArgs{
				ctx: context.Background(),
				in: &input.TodoDeleter{
					ID: todo.TodoID(2),
				},
			},
			wantErr: true,
		},
		"Internal error when soft deleting todo": {
			prepare: func(f *PrepareDeleteTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoHelper.
					EXPECT().
					CanAccessTodo(f.ctx, todo.UserID(1), todo.TodoID(1)).
					Return(true, nil).
					Times(1)

				f.mockTodoCommandsGateway.
					EXPECT().
					SoftDeleteTodo(f.ctx, todo.TodoID(1), todo.UserID(1)).
					Return(errors.New("database error")).
					Times(1)
			},
			args: DeleteTodoArgs{
				ctx: context.Background(),
				in: &input.TodoDeleter{
					ID: todo.TodoID(1),
				},
			},
			wantErr: true,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockBinder := mock_gateway.NewMockBinder(ctrl)
			mockTodoCommandsGateway := mock_gateway.NewMockTodoCommandsGateway(ctrl)
			mockTodoHelper := mock_service.NewMockTodoHelper(ctrl)

			todoDeleterService := service.NewTodoDeleter(
				mockBinder,
				mockTodoCommandsGateway,
				mockTodoHelper,
			)

			f := &PrepareDeleteTodoFields{
				ctx:                     context.Background(),
				mockBinder:              mockBinder,
				mockTodoCommandsGateway: mockTodoCommandsGateway,
				mockTodoHelper:          mockTodoHelper,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			err := todoDeleterService.SoftDelete(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("todoDeleterService.SoftDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
