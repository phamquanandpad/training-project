package auth

import (
	"google.golang.org/grpc"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	grpcutil "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/utils/grpc"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/config"
)

type TodoConnection grpc.ClientConnInterface

func NewAuthServiceClient(
	conn TodoConnection,
) auth_v1.AuthServiceClient {
	return auth_v1.NewAuthServiceClient(conn)
}

func NewAuthServiceConnection(conf *config.Config) (TodoConnection, func(), error) {
	return grpcutil.Dial(conf.AuthAddr, true)
}
