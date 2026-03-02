package service_test

import (
	"context"
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

type PrepareTodoGetterFields struct {
	ctx                    context.Context
	mockTodoQueriesGateway *mock_gateway.MockTodoQueriesGateway
}

type TodoGetterArgs struct {
	ctx context.Context
	in  *input.TodoGetter
}

type TodoGetterTestcase struct {
	prepare  func(f *PrepareTodoGetterFields)
	args     TodoGetterArgs
	expected *output.TodoGetter
	wantErr  bool
}

func Test_todoGetter_GetTodo(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testTables := map[string]TodoGetterTestcase{
		"Get Todo successfully": {
			prepare: func(f *PrepareTodoGetterFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					GetTodo(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.TodoID(1)).
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
			args: TodoGetterArgs{
				ctx: context.Background(),
				in: &input.TodoGetter{
					UserID: todo_model.UserID(1),
					ID:     todo_model.TodoID(1),
				},
			},
			expected: &output.TodoGetter{
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
		"Fail to get Todo when Todo does not exist": {
			prepare: func(f *PrepareTodoGetterFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					GetTodo(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, todo_model.TodoID(999)).
					Return(nil, nil).
					Times(1)
			},
			args: TodoGetterArgs{
				ctx: context.Background(),
				in: &input.TodoGetter{
					UserID: todo_model.UserID(1),
					ID:     todo_model.TodoID(999),
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

			mockTodoQueriesGateway := mock_gateway.NewMockTodoQueriesGateway(ctrl)

			f := &PrepareTodoGetterFields{
				ctx:                    context.Background(),
				mockTodoQueriesGateway: mockTodoQueriesGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewTodoGetter(mockTodoQueriesGateway)

			actual, err := svc.GetTodo(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("todoGetter.GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("todoGetter.GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			opts := cmpopts.IgnoreFields(todo_model.Todo{}, "CreatedAt", "UpdatedAt")
			if diff := cmp.Diff(tt.expected, actual, opts); diff != "" {
				t.Errorf("todoGetter.GetTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
