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

type PrepareTodoUpdaterFields struct {
	ctx                     context.Context
	mockTodoCommandsGateway *mock_gateway.MockTodoCommandsGateway
}

type TodoUpdaterArgs struct {
	ctx context.Context
	in  *input.TodoUpdater
}

type TodoUpdaterTestcase struct {
	prepare  func(f *PrepareTodoUpdaterFields)
	args     TodoUpdaterArgs
	expected *output.TodoUpdater
	wantErr  bool
}

func Test_todoUpdater_UpdateTodo(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testTables := map[string]TodoUpdaterTestcase{
		"Update Todo successfully": {
			prepare: func(f *PrepareTodoUpdaterFields) {
				status := todo_model.InProcess
				f.mockTodoCommandsGateway.
					EXPECT().
					Update(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.TodoID(1), todo_model.UpdateTodo{
						Task:        cast.Ptr("updated task"),
						Description: cast.Ptr("updated description"),
						Status:      &status,
					}).
					Return(&todo_model.Todo{
						ID:          todo_model.TodoID(1),
						UserID:      todo_model.UserID(1),
						Task:        "updated task",
						Description: cast.Ptr("updated description"),
						Status:      todo_model.InProcess,
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil).
					Times(1)
			},
			args: TodoUpdaterArgs{
				ctx: context.Background(),
				in: &input.TodoUpdater{
					UserID:      todo_model.UserID(1),
					ID:          todo_model.TodoID(1),
					Task:        cast.Ptr("updated task"),
					Description: cast.Ptr("updated description"),
					Status:      cast.Ptr(todo_model.InProcess),
				},
			},
			expected: &output.TodoUpdater{
				Todo: &todo_model.Todo{
					ID:          todo_model.TodoID(1),
					UserID:      todo_model.UserID(1),
					Task:        "updated task",
					Description: cast.Ptr("updated description"),
					Status:      todo_model.InProcess,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		"Fail to update Todo when Todo does not exist": {
			prepare: func(f *PrepareTodoUpdaterFields) {
				status := todo_model.InProcess
				f.mockTodoCommandsGateway.
					EXPECT().
					Update(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.TodoID(999), todo_model.UpdateTodo{
						Task:        cast.Ptr("updated task"),
						Description: cast.Ptr("updated description"),
						Status:      &status,
					}).
					Return(nil, nil).
					Times(1)
			},
			args: TodoUpdaterArgs{
				ctx: context.Background(),
				in: &input.TodoUpdater{
					UserID:      todo_model.UserID(1),
					ID:          todo_model.TodoID(999),
					Task:        cast.Ptr("updated task"),
					Description: cast.Ptr("updated description"),
					Status:      cast.Ptr(todo_model.InProcess),
				},
			},
			expected: nil,
			wantErr:  true,
		},
		"Fail to update Todo when gateway returns error": {
			prepare: func(f *PrepareTodoUpdaterFields) {
				status := todo_model.InProcess
				f.mockTodoCommandsGateway.
					EXPECT().
					Update(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.TodoID(1), todo_model.UpdateTodo{
						Task:        cast.Ptr("updated task"),
						Description: cast.Ptr("updated description"),
						Status:      &status,
					}).
					Return(nil, errors.New("gateway error")).
					Times(1)
			},
			args: TodoUpdaterArgs{
				ctx: context.Background(),
				in: &input.TodoUpdater{
					UserID:      todo_model.UserID(1),
					ID:          todo_model.TodoID(1),
					Task:        cast.Ptr("updated task"),
					Description: cast.Ptr("updated description"),
					Status:      cast.Ptr(todo_model.InProcess),
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

			f := &PrepareTodoUpdaterFields{
				ctx:                     context.Background(),
				mockTodoCommandsGateway: mockTodoCommandsGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewTodoUpdater(mockTodoCommandsGateway)

			actual, err := svc.UpdateTodo(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("todoUpdater.UpdateTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("todoUpdater.UpdateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			opts := cmpopts.IgnoreFields(todo_model.Todo{}, "CreatedAt", "UpdatedAt")
			if diff := cmp.Diff(tt.expected, actual, opts); diff != "" {
				t.Errorf("todoUpdater.UpdateTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
