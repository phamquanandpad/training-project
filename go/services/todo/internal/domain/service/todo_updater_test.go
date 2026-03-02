package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"
	mock_service "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type PrepareUpdateTodoFields struct {
	ctx                     context.Context
	mockBinder              *mock_gateway.MockBinder
	mockTodoCommandsGateway *mock_gateway.MockTodoCommandsGateway
	mockTodoHelper          *mock_service.MockTodoHelper
}

type UpdateTodoArgs struct {
	ctx context.Context
	in  *input.TodoUpdater
}

type UpdateTodoTestcase struct {
	name     string
	prepare  func(f *PrepareUpdateTodoFields)
	args     UpdateTodoArgs
	expected *output.TodoUpdater
	wantErr  bool
}

func Test_todoUpdater_Update(t *testing.T) {
	t.Parallel()

	createdAt := time.Now()
	existedUser := &todo.User{
		ID:       todo.UserID(1),
		Username: "user1",
	}

	testTables := map[string]UpdateTodoTestcase{
		"Update Todo successfully": {
			name: "Update Todo successfully",
			prepare: func(f *PrepareUpdateTodoFields) {
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
					UpdateTodo(f.ctx, todo.TodoID(1), todo.UserID(1), todo.UpdateTodo{
						Task:        cast.Ptr("updated task"),
						Description: cast.Ptr("updated description"),
						Status:      cast.Ptr(todo.Done),
					}).
					Return(&todo.Todo{
						ID:          todo.TodoID(1),
						UserID:      todo.UserID(1),
						Task:        "updated task",
						Description: cast.Ptr("updated description"),
						Status:      todo.Done,
						CreatedAt:   createdAt,
					}, nil).
					Times(1)
			},
			args: UpdateTodoArgs{
				ctx: context.Background(),
				in: &input.TodoUpdater{
					ID:          todo.TodoID(1),
					Task:        cast.Ptr("updated task"),
					Description: cast.Ptr("updated description"),
					Status:      cast.Ptr(todo.Done),
				},
			},
			expected: &output.TodoUpdater{
				ID:          todo.TodoID(1),
				UserID:      todo.UserID(1),
				Task:        "updated task",
				Description: cast.Ptr("updated description"),
				Status:      todo.Done,
				CreatedAt:   createdAt,
			},
			wantErr: false,
		},
		"Internal error from CanAccessTodo": {
			name: "Internal error from CanAccessTodo",
			prepare: func(f *PrepareUpdateTodoFields) {
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
			args: UpdateTodoArgs{
				ctx: context.Background(),
				in: &input.TodoUpdater{
					ID: todo.TodoID(1),
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Cannot access todo (unauthorized)": {
			name: "Cannot access todo (unauthorized)",
			prepare: func(f *PrepareUpdateTodoFields) {
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
			args: UpdateTodoArgs{
				ctx: context.Background(),
				in: &input.TodoUpdater{
					ID: todo.TodoID(2),
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Internal error when updating todo": {
			name: "Internal error when updating todo",
			prepare: func(f *PrepareUpdateTodoFields) {
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
					UpdateTodo(f.ctx, todo.TodoID(1), todo.UserID(1), gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: UpdateTodoArgs{
				ctx: context.Background(),
				in: &input.TodoUpdater{
					ID:   todo.TodoID(1),
					Task: cast.Ptr("updated task"),
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockBinder := mock_gateway.NewMockBinder(ctrl)
			mockTodoCommandsGateway := mock_gateway.NewMockTodoCommandsGateway(ctrl)
			mockTodoHelper := mock_service.NewMockTodoHelper(ctrl)

			todoUpdaterService := service.NewTodoUpdater(
				mockBinder,
				mockTodoCommandsGateway,
				mockTodoHelper,
			)

			f := &PrepareUpdateTodoFields{
				ctx:                     context.Background(),
				mockBinder:              mockBinder,
				mockTodoCommandsGateway: mockTodoCommandsGateway,
				mockTodoHelper:          mockTodoHelper,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := todoUpdaterService.Update(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("todoUpdaterService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(res, tt.expected); diff != "" {
				t.Errorf("UpdateTodo result mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
