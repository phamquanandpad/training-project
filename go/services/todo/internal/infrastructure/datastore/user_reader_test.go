package datastore_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/infrastructure/datastore"
)

func Test_userReader_GetUser(t *testing.T) {
	type args struct {
		userID todo.UserID
	}

	type testcase struct {
		args     args
		expected *todo.User
		wantErr  bool
	}

	t.Parallel()

	testTables := map[string]testcase{
		"Get User 1": {
			args: args{
				userID: 1,
			},
			expected: &todo.User{
				ID:        1,
				Username:  "user1",
				Email:     cast.Ptr("user1@example.com"),
				CreatedAt: getLocalTimeByString("2026-01-01T00:00:00Z"),
				UpdatedAt: getLocalTimeByString("2026-01-01T00:00:00Z"),
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
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userReader := datastore.NewUserReader()

			actual, err := userReader.GetUser(ctxWithReadDB, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userReader.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(actual, tt.expected); diff != "" {
				t.Fatalf("mismatch (-actual +expected):\n%s", diff)
			}
		})
	}
}
