package service_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type PrepareGetTodoFields struct {
	ctx                    context.Context
	mockBinder             *mock_gateway.MockBinder
	mockTodoQueriesGateway *mock_gateway.MockTodoQueriesGateway
}

type GetTodoArgs struct {
	ctx context.Context
	in  *input.TodoGetter
}

type GetTodoTestcase struct {
	prepare  func(f *PrepareGetTodoFields)
	args     GetTodoArgs
	expected *output.TodoGetter
	wantErr  bool
}

func Test_todoGetter_Get(t *testing.T) {
	t.Parallel()

	existedUser := &todo.User{
		ID:       todo.UserID(1),
		Username: "user1",
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	testTables := map[string]GetTodoTestcase{
		"Get Todo successfully": {
			prepare: func(f *PrepareGetTodoFields) {
				f.ctx = todo.WithUser(f.ctx, existedUser)
				f.mockBinder.
					EXPECT().
					Bind(gomock.Any()).
					Return(f.ctx).
					Times(1)

				f.mockTodoQueriesGateway.
					EXPECT().
					GetTodo(f.ctx, todo.TodoID(1), todo.UserID(1)).
					Return(&todo.Todo{
						ID:          1,
						UserID:      todo.UserID(1),
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo.Pending,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					}, nil)
			},
			args: GetTodoArgs{
				ctx: context.Background(),
				in: &input.TodoGetter{
					ID: 1,
				},
			},
			expected: &output.TodoGetter{
				ID:          1,
				Task:        "todo task 1",
				Description: cast.Ptr("todo description 1"),
				Status:      todo.Pending,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			mockBinder := mock_gateway.NewMockBinder(ctrl)
			mockTodoQueriesGateway := mock_gateway.NewMockTodoQueriesGateway(ctrl)
			todoGetterService := service.NewTodoGetter(mockBinder, mockTodoQueriesGateway)
			f := &PrepareGetTodoFields{
				ctx:                    context.Background(),
				mockBinder:             mockBinder,
				mockTodoQueriesGateway: mockTodoQueriesGateway,
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}

			realResult, err := todoGetterService.Get(f.ctx, tt.args.in)
			if !reflect.DeepEqual(realResult, tt.expected) {
				t.Errorf("GetReport() realResult = %v is not equal to expected %v", realResult, tt.expected)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetReport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
