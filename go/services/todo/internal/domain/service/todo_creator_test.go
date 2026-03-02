package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type PrepareCreateTodoFields struct {
	ctx                     context.Context
	mockBinder              *mock_gateway.MockBinder
	mockTodoCommandsGateway *mock_gateway.MockTodoCommandsGateway
}

type CreateTodoArgs struct {
	ctx context.Context
	in  *input.TodoCreator
}

type CreateTodoTestcase struct {
	name     string
	prepare  func(f *PrepareCreateTodoFields)
	args     CreateTodoArgs
	expected *output.TodoCreator
	wantErr  bool
}

func Test_todoCreator_Create(t *testing.T) {
	t.Parallel()

	existedUser := &todo.User{
		ID:       todo.UserID(1),
		Username: "user1",
	}

	testTables := map[string]CreateTodoTestcase{
		"Create Todo successfully": {
			name: "Create Todo successfully",
			prepare: func(f *PrepareCreateTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoCommandsGateway.
					EXPECT().
					CreateTodo(f.ctx, todo.NewTodo{
						UserID:      todo.UserID(1),
						Task:        "new task",
						Description: cast.Ptr("new description"),
						Status:      todo.Pending,
					}).
					Return(&todo.Todo{
						UserID:      todo.UserID(1),
						Task:        "new task",
						Description: cast.Ptr("new description"),
						Status:      todo.Pending,
					}, nil).
					Times(1)
			},
			args: CreateTodoArgs{
				ctx: context.Background(),
				in: &input.TodoCreator{
					Task:        "new task",
					Description: cast.Ptr("new description"),
					Status:      todo.Pending,
				},
			},
			expected: &output.TodoCreator{
				UserID:      todo.UserID(1),
				Task:        "new task",
				Description: cast.Ptr("new description"),
				Status:      todo.Pending,
			},
			wantErr: false,
		},
		"Create Todo without description": {
			name: "Create Todo without description",
			prepare: func(f *PrepareCreateTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoCommandsGateway.
					EXPECT().
					CreateTodo(f.ctx, todo.NewTodo{
						UserID:      todo.UserID(1),
						Task:        "another task",
						Description: nil,
						Status:      todo.InProcess,
					}).
					Return(&todo.Todo{
						UserID:      todo.UserID(1),
						Task:        "another task",
						Description: nil,
						Status:      todo.InProcess,
					}, nil).
					Times(1)
			},
			args: CreateTodoArgs{
				ctx: context.Background(),
				in: &input.TodoCreator{
					Task:   "another task",
					Status: todo.InProcess,
				},
			},
			expected: &output.TodoCreator{
				UserID:      todo.UserID(1),
				Task:        "another task",
				Description: nil,
				Status:      todo.InProcess,
			},
			wantErr: false,
		},
		"Internal error when creating todo": {
			name: "Internal error when creating todo",
			prepare: func(f *PrepareCreateTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoCommandsGateway.
					EXPECT().
					CreateTodo(f.ctx, gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: CreateTodoArgs{
				ctx: context.Background(),
				in: &input.TodoCreator{
					Task:   "new task",
					Status: todo.Pending,
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

			todoCreatorService := service.NewTodoCreator(
				mockBinder,
				mockTodoCommandsGateway,
			)

			f := &PrepareCreateTodoFields{
				ctx:                     context.Background(),
				mockBinder:              mockBinder,
				mockTodoCommandsGateway: mockTodoCommandsGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := todoCreatorService.Create(f.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("todoCreatorService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(todo.Todo{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt"),
			}

			if diff := cmp.Diff(res, tt.expected, ignoreFieldsOpts...); diff != "" {
				t.Errorf("CreateTodo result mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
