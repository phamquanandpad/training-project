package datastore_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/infrastructure/datastore"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/testutil"
)

func Test_userWriter_CreateUser(t *testing.T) {
	t.Parallel()
	gormDB, _ := testutil.InitDB(t)

	type args struct {
		id        todo.UserID
		username  string
		email     *string
		createdAt time.Time
		updatedAt time.Time
	}

	type testcase struct {
		args     args
		expected *todo.User
		wantErr  bool
	}

	createdAt := time.Now()
	updatedAt := time.Now()
	testTables := map[string]testcase{
		"Create User": {
			args: args{
				id:        10,
				username:  "newuser10",
				email:     cast.Ptr("newuser10@example.com"),
				createdAt: createdAt,
				updatedAt: updatedAt,
			},
			expected: &todo.User{
				ID:        todo.UserID(10),
				Username:  "newuser10",
				Email:     cast.Ptr("newuser10@example.com"),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userWriter := datastore.NewUserWriter()

			ctx := datastore.WithTodoDB(context.Background(), gormDB)
			actual, err := userWriter.CreateUser(ctx, todo.NewUser{
				ID:        tt.args.id,
				Username:  tt.args.username,
				Email:     tt.args.email,
				CreatedAt: tt.args.createdAt,
				UpdatedAt: tt.args.updatedAt,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("userWriter.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(todo.User{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt"),
			}

			if diff := cmp.Diff(tt.expected, actual, ignoreFieldsOpts...); diff != "" {
				t.Errorf("userWriter.CreateUser() = diff: %s", diff)
			}
		})
	}
}

func Test_userWriter_CreateUser_Duplicate(t *testing.T) {
	t.Parallel()
	gormDB, _ := testutil.InitDB(t)

	type args struct {
		id        todo.UserID
		username  string
		email     *string
		createdAt time.Time
		updatedAt time.Time
	}

	type testcase struct {
		args     args
		expected *todo.User
		wantErr  bool
	}

	createdAt := time.Now()
	updatedAt := time.Now()
	testTables := map[string]testcase{
		"Create User with duplicate email": {
			args: args{
				id:        todo.UserID(10),
				username:  "newuser10",
				email:     cast.Ptr("user1@example.com"),
				createdAt: createdAt,
				updatedAt: updatedAt,
			},
			expected: nil,
			wantErr:  true,
		},
		"Create User with duplicate username": {
			args: args{
				id:        todo.UserID(11),
				username:  "user1",
				email:     cast.Ptr("newuser11@example.com"),
				createdAt: createdAt,
				updatedAt: updatedAt,
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for name, tt := range testTables {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userWriter := datastore.NewUserWriter()

			ctx := datastore.WithTodoDB(context.Background(), gormDB)
			actual, err := userWriter.CreateUser(ctx, todo.NewUser{
				ID:        tt.args.id,
				Username:  tt.args.username,
				Email:     tt.args.email,
				CreatedAt: tt.args.createdAt,
				UpdatedAt: tt.args.updatedAt,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("userWriter.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t.Errorf("userWriter.CreateUser() expected error but got nil")
				return
			}
			if actual != nil {
				t.Errorf("userWriter.CreateUser() expected nil but got %v", actual)
			}
		})
	}
}
