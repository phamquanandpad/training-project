package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	mock_gateway "github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway/mock"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
)

func Test_todoHelper_CanAccessTodo(t *testing.T) {
	type prepareFields struct {
		ctx                    context.Context
		mockTodoQueriesGateway *mock_gateway.MockTodoQueriesGateway
	}

	type args struct {
		ctx    context.Context
		userID todo.UserID
		todoID todo.TodoID
	}

	type testcase struct {
		prepare  func(f *prepareFields)
		args     args
		expected bool
		wantErr  bool
	}

	testTables := map[string]testcase{
		"Can access todo": {
			prepare: func(f *prepareFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					GetTodo(f.ctx, todo.TodoID(1), todo.UserID(1)).
					Return(&todo.Todo{
						ID:     todo.TodoID(1),
						UserID: todo.UserID(1),
					}, nil)
			},
			args: args{
				ctx:    context.Background(),
				userID: todo.UserID(1),
				todoID: todo.TodoID(1),
			},
			expected: true,
			wantErr:  false,
		},
		"Cannot access todo due to different user": {
			prepare: func(f *prepareFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					GetTodo(f.ctx, todo.TodoID(2), todo.UserID(2)).
					Return(nil, nil)
			},
			args: args{
				ctx:    context.Background(),
				userID: todo.UserID(2),
				todoID: todo.TodoID(2),
			},
			expected: false,
			wantErr:  true,
		},
		"Todo not found": {
			prepare: func(f *prepareFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					GetTodo(f.ctx, todo.TodoID(999), todo.UserID(1)).
					Return(nil, nil)
			},
			args: args{
				ctx:    context.Background(),
				userID: todo.UserID(1),
				todoID: todo.TodoID(999),
			},
			expected: false,
			wantErr:  true,
		},
		"Internal error when getting todo": {
			prepare: func(f *prepareFields) {
				f.mockTodoQueriesGateway.
					EXPECT().
					GetTodo(f.ctx, todo.TodoID(1000), todo.UserID(1)).
					Return(nil, errors.New("database error"))
			},
			args: args{
				ctx:    context.Background(),
				userID: todo.UserID(1),
				todoID: todo.TodoID(1000),
			},
			expected: false,
			wantErr:  true,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			t.Cleanup(ctrl.Finish)

			mockTodoQueriesGateway := mock_gateway.NewMockTodoQueriesGateway(ctrl)

			todoHelper := service.NewTodoHelper(
				mockTodoQueriesGateway,
			)

			f := &prepareFields{
				ctx:                    tt.args.ctx,
				mockTodoQueriesGateway: mockTodoQueriesGateway,
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			actual, err := todoHelper.CanAccessTodo(
				tt.args.ctx,
				tt.args.userID,
				tt.args.todoID,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("todoHelper.CanAccessTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("todoHelper.CanAccessTodo() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}
