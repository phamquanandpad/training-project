package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway/mock"
	todo_model "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type PrepareTodoCreatorFields struct {
	ctx                     context.Context
	mockTodoCommandsGateway *mock_gateway.MockTodoCommandsGateway
}

type TodoCreatorArgs struct {
	ctx context.Context
	in  *input.TodoCreator
}

type TodoCreatorTestcase struct {
	prepare  func(f *PrepareTodoCreatorFields)
	args     TodoCreatorArgs
	expected *output.TodoCreator
	wantErr  bool
}

func Test_todoCreator_Create(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testTables := map[string]TodoCreatorTestcase{
		"Create Todo successfully": {
			prepare: func(f *PrepareTodoCreatorFields) {
				f.mockTodoCommandsGateway.
					EXPECT().
					Create(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.NewTodo{
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo_model.Pending,
					}).
					Return(&todo_model.Todo{
						ID:          todo_model.TodoID(1),
						UserID:      todo_model.UserID(1),
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo_model.Pending,
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil).
					Times(1)
			},
			args: TodoCreatorArgs{
				ctx: context.Background(),
				in: &input.TodoCreator{
					UserID:      todo_model.UserID(1),
					Task:        "todo task 1",
					Description: cast.Ptr("todo description 1"),
					Status:      todo_model.Pending,
				},
			},
			expected: &output.TodoCreator{
				Todo: &todo_model.Todo{
					ID:          todo_model.TodoID(1),
					UserID:      todo_model.UserID(1),
					Task:        "todo task 1",
					Description: cast.Ptr("todo description 1"),
					Status:      todo_model.Pending,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		"Fail to create Todo when gateway returns error": {
			prepare: func(f *PrepareTodoCreatorFields) {
				f.mockTodoCommandsGateway.
					EXPECT().
					Create(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.NewTodo{
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo_model.Pending,
					}).
					Return(nil, errors.New("gateway error")).
					Times(1)
			},
			args: TodoCreatorArgs{
				ctx: context.Background(),
				in: &input.TodoCreator{
					UserID:      todo_model.UserID(1),
					Task:        "todo task 1",
					Description: cast.Ptr("todo description 1"),
					Status:      todo_model.Pending,
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

			mockTodoCommandsGateway := mock_gateway.NewMockTodoCommandsGateway(ctrl)

			f := &PrepareTodoCreatorFields{
				ctx:                     context.Background(),
				mockTodoCommandsGateway: mockTodoCommandsGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewTodoCreator(mockTodoCommandsGateway)

			actual, err := svc.CreateTodo(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("todoCreator.CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("todoCreator.CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			opts := cmpopts.IgnoreFields(todo_model.Todo{}, "CreatedAt", "UpdatedAt")
			if diff := cmp.Diff(tt.expected, actual, opts); diff != "" {
				t.Errorf("todoCreator.CreateTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
