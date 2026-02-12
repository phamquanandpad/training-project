package datastore_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/infrastructure/datastore"
)

func Test_todoReader_GetTodo(t *testing.T) {
	type args struct {
		todoID todo.TodoID
		userID todo.UserID
	}

	type testcase struct {
		args     args
		expected *todo.Todo
		wantErr  bool
	}

	t.Parallel()

	testTables := map[string]testcase{
		"Get Todo 1": {
			args: args{
				todoID: 1,
				userID: 1,
			},
			expected: &todo.Todo{
				ID:          1,
				UserID:      todo.UserID(1),
				Task:        "todo task 1",
				Description: cast.Ptr("todo description 1"),
				Status:      todo.Pending, // 0
				CreatedAt:   getLocalTimeByString("2026-01-01T00:00:00Z"),
				UpdatedAt:   getLocalTimeByString("2026-01-01T00:00:00Z"),
			},
			wantErr: false,
		},
		"Not found and return nil": {
			args: args{
				todoID: 999,
				userID: 1,
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for name, tt := range testTables {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			todoReader := datastore.NewTodoReader()

			actual, err := todoReader.GetTodo(ctxWithReadDB, tt.args.todoID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(actual, tt.expected); diff != "" {
				t.Fatalf("mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}

func Test_todoReader_ListTodos(t *testing.T) {
	type args struct {
		userID todo.UserID
	}

	type expected struct {
		todos []*todo.Todo
		total int
	}

	type testcase struct {
		args     args
		expected expected
		wantErr  bool
	}

	t.Parallel()

	testTables := map[string]testcase{
		"List Todos for User 1": {
			args: args{userID: 1},
			expected: expected{
				todos: []*todo.Todo{
					{
						ID:          2,
						UserID:      1,
						Task:        "todo task 2",
						Description: cast.Ptr("todo description 2"),
						Status:      todo.InProcess,
						CreatedAt:   getLocalTimeByString("2026-01-02T00:00:00Z"),
						UpdatedAt:   getLocalTimeByString("2026-01-02T00:00:00Z"),
					},
					{
						ID:          1,
						UserID:      1,
						Task:        "todo task 1",
						Description: cast.Ptr("todo description 1"),
						Status:      todo.Pending,
						CreatedAt:   getLocalTimeByString("2026-01-01T00:00:00Z"),
						UpdatedAt:   getLocalTimeByString("2026-01-01T00:00:00Z"),
					},
				},
				total: 2,
			},
			wantErr: false,
		},
		"List Todos for User 2": {
			args: args{userID: 2},
			expected: expected{
				todos: []*todo.Todo{
					{
						ID:          3,
						UserID:      2,
						Task:        "todo task 3",
						Description: cast.Ptr("todo description 3"),
						Status:      todo.Pending,
						CreatedAt:   getLocalTimeByString("2026-01-03T00:00:00Z"),
						UpdatedAt:   getLocalTimeByString("2026-01-03T00:00:00Z"),
					},
				},
				total: 1,
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			todoReader := datastore.NewTodoReader()

			todos, total, err := todoReader.ListTodos(ctxWithReadDB, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(todos, tt.expected.todos); diff != "" {
				t.Fatalf("todos mismatch (-actual +expected):\n%s", diff)
			}

			if total != tt.expected.total {
				t.Fatalf("total = %d want %d", total, tt.expected.total)
			}
		})
	}
}
