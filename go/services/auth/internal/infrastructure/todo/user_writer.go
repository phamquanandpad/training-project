package todo

import (
	"context"

	todo_common_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/common/v1"
	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/todo"
)

type userWriter struct {
	todoServiceClient todo_v1.TodoServiceClient
}

func NewUserWriter(
	todoServiceClient todo_v1.TodoServiceClient,
) gateway.UserServiceCommandsGateway {
	return &userWriter{
		todoServiceClient: todoServiceClient,
	}
}

func (w userWriter) CreateUser(
	ctx context.Context,
	newUser todo.NewUser,
) (*todo.User, error) {
	_, err := w.todoServiceClient.PostUser(
		ctx,
		&todo_v1.PostUserRequest{
			User: &todo_common_v1.User{
				Id:       newUser.ID.Int64(),
				Username: newUser.Username,
				Email:    cast.Value(newUser.Email),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &todo.User{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
	}, nil
}
