package datastore_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/datastore"
)

func Test_userReader_GetUserByID(t *testing.T) {
	type args struct {
		userID auth_models.UserID
	}

	type testcase struct {
		args     args
		expected *auth_models.User
		wantErr  bool
	}

	t.Parallel()

	testTables := map[string]testcase{
		"Get User 1 by ID": {
			args: args{
				userID: 1,
			},
			expected: &auth_models.User{
				ID:        1,
				Email:     "user1@example.com",
				CreatedAt: getLocalTimeByString("2026-01-01T00:00:00Z"),
				UpdatedAt: getLocalTimeByString("2026-01-01T00:00:00Z"),
			},
			wantErr: false,
		},
		"Get User 2 by ID": {
			args: args{
				userID: 2,
			},
			expected: &auth_models.User{
				ID:        2,
				Email:     "user2@example.com",
				CreatedAt: getLocalTimeByString("2026-01-02T00:00:00Z"),
				UpdatedAt: getLocalTimeByString("2026-01-02T00:00:00Z"),
			},
			wantErr: false,
		},
		"Not found and return nil": {
			args: args{
				userID: 999,
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userReader := datastore.NewUserReader()

			actual, err := userReader.GetUserByID(ctxWithReadDB, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userReader.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ignorePasswordOpt := cmp.FilterPath(func(fp cmp.Path) bool {
				return fp.Last().String() == ".Password"
			}, cmp.Ignore())

			if diff := cmp.Diff(actual, tt.expected, ignorePasswordOpt); diff != "" {
				t.Fatalf("mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}

func Test_userReader_GetUserByEmail(t *testing.T) {
	type args struct {
		email string
	}

	type testcase struct {
		args     args
		expected *auth_models.User
		wantErr  bool
	}

	t.Parallel()

	testTables := map[string]testcase{
		"Get User by email user1@example.com": {
			args: args{
				email: "user1@example.com",
			},
			expected: &auth_models.User{
				ID:        1,
				Email:     "user1@example.com",
				CreatedAt: getLocalTimeByString("2026-01-01T00:00:00Z"),
				UpdatedAt: getLocalTimeByString("2026-01-01T00:00:00Z"),
			},
			wantErr: false,
		},
		"Get User by email user2@example.com": {
			args: args{
				email: "user2@example.com",
			},
			expected: &auth_models.User{
				ID:        2,
				Email:     "user2@example.com",
				CreatedAt: getLocalTimeByString("2026-01-02T00:00:00Z"),
				UpdatedAt: getLocalTimeByString("2026-01-02T00:00:00Z"),
			},
			wantErr: false,
		},
		"Not found and return nil": {
			args: args{
				email: "notfound@example.com",
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userReader := datastore.NewUserReader()

			actual, err := userReader.GetUserByEmail(ctxWithReadDB, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userReader.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ignorePasswordOpt := cmp.FilterPath(func(fp cmp.Path) bool {
				return fp.Last().String() == ".Password"
			}, cmp.Ignore())

			if diff := cmp.Diff(actual, tt.expected, ignorePasswordOpt); diff != "" {
				t.Fatalf("mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}
