package mapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/mapper"
)

func Test_ToUserGRPCResponse(t *testing.T) {
	type args struct {
		user *todo.User
	}

	testTables := map[string]struct {
		args     args
		expected *todo_common_v1.User
	}{
		"Convert User to gRPC response successfully": {
			args: args{
				user: &todo.User{
					ID:       todo.UserID(1),
					Username: "testuser",
					Email:    cast.Ptr("test@example.com"),
				},
			},
			expected: &todo_common_v1.User{
				Id:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
		"Convert User with nil Email to gRPC response with empty Email": {
			args: args{
				user: &todo.User{
					ID:       todo.UserID(2),
					Username: "anotheruser",
					Email:    nil,
				},
			},
			expected: &todo_common_v1.User{
				Id:       2,
				Username: "anotheruser",
				Email:    "",
			},
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			out := mapper.ToUserGRPCResponse(tt.args.user)
			if diff := cmp.Diff(out, tt.expected, protocmp.Transform()); diff != "" {
				t.Errorf("ToUserGRPCResponse() = got differs from expected: %s", diff)
			}
		})
	}
}
