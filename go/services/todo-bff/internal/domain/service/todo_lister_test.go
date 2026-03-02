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

type PrepareTodoListerFields struct {
	ctx                    context.Context
	mockTodoQueriesGateway *mock_gateway.MockTodoQueriesGateway
}

type TodoListerArgs struct {
	ctx context.Context
	in  *input.TodoLister
}

type TodoListerTestcase struct {
	prepare  func(f *PrepareTodoListerFields)
	args     TodoListerArgs
	expected *output.TodoLister
	wantErr  bool
}

func Test_todoLister_ListTodos(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testTables := map[string]TodoListerTestcase{
		"List Todos successfully": {
			prepare: func(f *PrepareTodoListerFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					ListTodos(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, 10, 0).
					Return([]*todo_model.Todo{
						{
							ID:          todo_model.TodoID(1),
							UserID:      todo_model.UserID(1),
							Task:        "todo task 1",
							Description: cast.Ptr("todo description 1"),
							Status:      todo_model.Pending,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
						{
							ID:          todo_model.TodoID(2),
							UserID:      todo_model.UserID(1),
							Task:        "todo task 2",
							Description: nil,
							Status:      todo_model.InProcess,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					}, 2, nil).
					Times(1)
			},
			args: TodoListerArgs{
				ctx: context.Background(),
				in: &input.TodoLister{
					UserID: todo_model.UserID(1),
					Limit:  10,
					Offset: 0,
				},
			},
			expected: &output.TodoLister{
				Todos: []*todo_model.Todo{
					{
						ID:          todo_model.TodoID(1),
						UserID:      todo_model.UserID(1),
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo_model.Pending,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						ID:          todo_model.TodoID(2),
						UserID:      todo_model.UserID(1),
						Task:        "todo task 2",
						Description: nil,
						Status:      todo_model.InProcess,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
				TotalCount: 2,
			},
			wantErr: false,
		},
		"List Todos successfully with empty list": {
			prepare: func(f *PrepareTodoListerFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					ListTodos(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, 10, 0).
					Return([]*todo_model.Todo{}, 0, nil).
					Times(1)
			},
			args: TodoListerArgs{
				ctx: context.Background(),
				in: &input.TodoLister{
					UserID: todo_model.UserID(1),
					Limit:  10,
					Offset: 0,
				},
			},
			expected: &output.TodoLister{
				Todos:      []*todo_model.Todo{},
				TotalCount: 0,
			},
			wantErr: false,
		},
		"Fail to list Todos when gateway returns error": {
			prepare: func(f *PrepareTodoListerFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					ListTodos(f.ctx, todo_model.UserAttributes{
						UserID: todo_model.UserID(1),
					}, 10, 0).
					Return(nil, 0, errors.New("gateway error")).
					Times(1)
			},
			args: TodoListerArgs{
				ctx: context.Background(),
				in: &input.TodoLister{
					UserID: todo_model.UserID(1),
					Limit:  10,
					Offset: 0,
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

			f := &PrepareTodoListerFields{
				ctx:                    context.Background(),
				mockTodoQueriesGateway: mockTodoQueriesGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			svc := service.NewTodoLister(mockTodoQueriesGateway)

			actual, err := svc.ListTodos(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("todoLister.ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("todoLister.ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			opts := cmpopts.IgnoreFields(todo_model.Todo{}, "CreatedAt", "UpdatedAt")
			if diff := cmp.Diff(tt.expected, actual, opts); diff != "" {
				t.Errorf("todoLister.ListTodos() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
