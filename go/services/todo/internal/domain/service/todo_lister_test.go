package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type PrepareListTodoFields struct {
	ctx                    context.Context
	mockBinder             *mock_gateway.MockBinder
	mockTodoQueriesGateway *mock_gateway.MockTodoQueriesGateway
}

type ListTodoArgs struct {
	ctx context.Context
	in  *input.TodoLister
}

type ListTodoTestcase struct {
	name     string
	prepare  func(f *PrepareListTodoFields)
	args     ListTodoArgs
	expected *output.TodoLister
	wantErr  bool
}

func Test_todoLister_List(t *testing.T) {
	t.Parallel()
	now := time.Now()
	createdAt := now
	updatedAt := now
	limit := int64(10)
	offset := int64(0)
	testTables := []ListTodoTestcase{
		{
			name: "List Todos successfully",
			prepare: func(f *PrepareListTodoFields) {
				existedUser := &todo.User{
					ID: todo.UserID(1),
				}
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoQueriesGateway.
					EXPECT().
					ListTodos(f.ctx, todo.UserID(1), int(limit), int(offset)).
					Return([]*todo.Todo{
						{
							ID:          todo.TodoID(1),
							UserID:      todo.UserID(1),
							Task:        "todo task 1",
							Description: cast.Ptr("todo description 1"),
							Status:      todo.Pending,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
						{
							ID:          todo.TodoID(2),
							UserID:      todo.UserID(1),
							Task:        "todo task 2",
							Description: nil,
							Status:      todo.InProcess,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
					}, 2, nil).
					Times(1)
			},
			args: ListTodoArgs{
				ctx: context.Background(),
				in: &input.TodoLister{
					UserAttributes: input.UserAttributes{
						UserID: todo.UserID(1),
					},
					Limit:  limit,
					Offset: offset,
				},
			},

			expected: &output.TodoLister{
				Todos: []*todo.Todo{
					{
						ID:          todo.TodoID(1),
						UserID:      todo.UserID(1),
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo.Pending,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
					{
						ID:          todo.TodoID(2),
						UserID:      todo.UserID(1),
						Task:        "todo task 2",
						Description: nil,
						Status:      todo.InProcess,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
				},
				Total: 2,
			},
			wantErr: false,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			t.Cleanup(ctrl.Finish)

			mockBinder := mock_gateway.NewMockBinder(ctrl)
			mockTodoQueriesGateway := mock_gateway.NewMockTodoQueriesGateway(ctrl)

			todoListerService := service.NewTodoLister(
				mockBinder,
				mockTodoQueriesGateway,
			)

			f := &PrepareListTodoFields{
				ctx:                    context.Background(),
				mockBinder:             mockBinder,
				mockTodoQueriesGateway: mockTodoQueriesGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			res, err := todoListerService.List(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("todoListerService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(res, tt.expected); diff != "" {
				t.Errorf("ListTodos result mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}
