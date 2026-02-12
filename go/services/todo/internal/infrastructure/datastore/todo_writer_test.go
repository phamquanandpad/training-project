package datastore_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/infrastructure/datastore"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/testutil"
)

func Test_todoWriter_CreateTodo(t *testing.T) {
	t.Parallel()
	gormDB, _ := testutil.InitDB(t)

	type args struct {
		userID      todo.UserID
		task        string
		description *string
		status      todo.TodoStatus
	}

	type testcase struct {
		args     args
		expected *todo.Todo
		wantErr  bool
	}

	testTables := map[string]testcase{
		"Create Todo": {
			args: args{
				userID:      todo.UserID(1),
				task:        "new todo task 1",
				description: cast.Ptr("new todo description 1"),
				status:      todo.Pending,
			},
			expected: &todo.Todo{
				UserID:      todo.UserID(1),
				Task:        "new todo task 1",
				Description: cast.Ptr("new todo description 1"),
				Status:      todo.Pending,
			},
			wantErr: false,
		},
		"Create Todo with no description": {
			args: args{
				userID:      todo.UserID(1),
				task:        "new todo task 2",
				description: nil,
				status:      todo.Pending,
			},
			expected: &todo.Todo{
				UserID:      todo.UserID(1),
				Task:        "new todo task 2",
				Description: nil,
				Status:      todo.Pending,
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			tx := gormDB.Begin()

			defer tx.Rollback()

			ctxWithWriteDB := datastore.WithTodoDB(context.Background(), tx)

			todoWriter := datastore.NewTodoWriter()
			res, err := todoWriter.CreateTodo(ctxWithWriteDB, todo.NewTodo{
				UserID:      tt.args.userID,
				Task:        tt.args.task,
				Description: tt.args.description,
				Status:      tt.args.status,
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: got %v, wantErr %v", err, tt.wantErr)
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(todo.Todo{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt"),
			}

			if diff := cmp.Diff(tt.expected, res, ignoreFieldsOpts...); diff != "" {
				t.Errorf("todoWriter.CreateTodo() value is mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}

func Test_todoWriter_UpdateTodo(t *testing.T) {
	t.Parallel()
	gormDB, _ := testutil.InitDB(t)

	type args struct {
		todoID     todo.TodoID
		userID     todo.UserID
		updateTodo todo.UpdateTodo
	}

	type testcase struct {
		args     args
		expected *todo.Todo
		wantErr  bool
	}

	testTables := map[string]testcase{
		"Update Todo by User return success": {
			args: args{
				todoID: todo.TodoID(1),
				userID: todo.UserID(1),
				updateTodo: todo.UpdateTodo{
					Task: cast.Ptr("updated todo task 1"),
				},
			},
			expected: &todo.Todo{
				ID:          todo.TodoID(1),
				UserID:      todo.UserID(1),
				Task:        "updated todo task 1",
				Description: cast.Ptr("todo description 1"),
				Status:      todo.Pending,
			},
			wantErr: false,
		},
		"Update Todo by User return error when todo is not created by the User": {
			args: args{
				todoID: todo.TodoID(3),
				userID: todo.UserID(1),
				updateTodo: todo.UpdateTodo{
					Description: cast.Ptr("updated todo description 2"),
				},
			},
			expected: nil,
			wantErr:  false,
		},
		"Update Todo by User return nil when todo not found": {
			args: args{
				todoID: todo.TodoID(999),
				userID: todo.UserID(1),
				updateTodo: todo.UpdateTodo{
					Status: cast.Ptr(todo.Done),
				},
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for name, tt := range testTables {
		tt := tt
		t.Run(name, func(t *testing.T) {
			tx := gormDB.Begin()

			defer tx.Rollback()

			ctxWithWriteDB := datastore.WithTodoDB(context.Background(), tx)
			todoWriter := datastore.NewTodoWriter()
			res, err := todoWriter.UpdateTodo(
				ctxWithWriteDB,
				tt.args.todoID,
				tt.args.userID,
				tt.args.updateTodo,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("todoWriter.UpdateTodo() error = %v, wantErr %v", err, tt.wantErr)
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(todo.Todo{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
			}

			if diff := cmp.Diff(tt.expected, res, ignoreFieldsOpts...); diff != "" {
				t.Errorf("todoWriter.UpdateTodo() value is mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}

func Test_todoWriter_SoftDeleteTodo(t *testing.T) {
	t.Parallel()
	gormDB, _ := testutil.InitDB(t)

	type args struct {
		todoID todo.TodoID
		userID todo.UserID
	}

	type testcase struct {
		args    args
		wantErr bool
	}

	testTables := map[string]testcase{
		"Soft Delete Todo by User return success": {
			args: args{
				todoID: todo.TodoID(2),
				userID: todo.UserID(1),
			},
			wantErr: false,
		},
		"Soft Delete Todo by User return error when todo is not created by the User": {
			args: args{
				todoID: todo.TodoID(3),
				userID: todo.UserID(1),
			},
			wantErr: false,
		},
		"Soft Delete Todo by User return nil when todo not found": {
			args: args{
				todoID: todo.TodoID(999),
				userID: todo.UserID(1),
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		tt := tt
		t.Run(name, func(t *testing.T) {
			tx := gormDB.Begin()

			defer tx.Rollback()

			ctxWithWriteDB := datastore.WithTodoDB(context.Background(), tx)
			todoWriter := datastore.NewTodoWriter()
			err := todoWriter.SoftDeleteTodo(
				ctxWithWriteDB,
				tt.args.todoID,
				tt.args.userID,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("todoWriter.SoftDeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
