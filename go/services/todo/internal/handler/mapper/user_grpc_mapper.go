package mapper

import (
	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

func ToUserGRPCResponse(user *todo.User) *todo_common_v1.User {
	return &todo_common_v1.User{
		Id:       int64(user.ID),
		Username: user.Username,
		Email:    cast.Value(user.Email),
	}
}
