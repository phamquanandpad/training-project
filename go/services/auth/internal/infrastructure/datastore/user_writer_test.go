package datastore_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/datastore"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/testutil"
)

func Test_userWriter_CreateUser(t *testing.T) {
	t.Parallel()
	gormDB, _ := testutil.InitDB(t)

	type args struct {
		email    string
		password string
	}

	type testcase struct {
		args     args
		expected *auth_models.User
		wantErr  bool
	}

	testTables := map[string]testcase{
		"Create User successfully": {
			args: args{
				email:    "newuser10@example.com",
				password: "hashedpassword",
			},
			expected: &auth_models.User{
				Email:    "newuser10@example.com",
				Password: "hashedpassword",
			},
			wantErr: false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userWriter := datastore.NewUserWriter()

			ctx := datastore.WithAuthDB(context.Background(), gormDB)
			actual, err := userWriter.CreateUser(ctx, auth_models.NewUser{
				Email:    tt.args.email,
				Password: tt.args.password,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("userWriter.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			ignoreFieldsOpts := []cmp.Option{
				cmpopts.IgnoreFields(auth_models.User{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt"),
			}

			if diff := cmp.Diff(tt.expected, actual, ignoreFieldsOpts...); diff != "" {
				t.Errorf("userWriter.CreateUser() diff: %s", diff)
			}
		})
	}
}

func Test_userWriter_CreateUser_Duplicate(t *testing.T) {
	t.Parallel()
	gormDB, _ := testutil.InitDB(t)

	type args struct {
		email    string
		password string
	}

	type testcase struct {
		args    args
		wantErr bool
	}

	testTables := map[string]testcase{
		"Create User with duplicate email": {
			args: args{
				email:    "user1@example.com",
				password: "hashedpassword",
			},
			wantErr: true,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userWriter := datastore.NewUserWriter()

			ctx := datastore.WithAuthDB(context.Background(), gormDB)
			_, err := userWriter.CreateUser(ctx, auth_models.NewUser{
				Email:    tt.args.email,
				Password: tt.args.password,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("userWriter.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
