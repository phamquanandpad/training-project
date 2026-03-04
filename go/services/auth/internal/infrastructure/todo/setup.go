package todo

import (
	"google.golang.org/grpc"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	grpcutil "github.com/phamquanandpad/training-project/go/services/auth/internal/utils/grpc"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/config"
)

type TodoConnection grpc.ClientConnInterface

func NewTodoServiceClient(
	conn TodoConnection,
) todo_v1.TodoServiceClient {
	return todo_v1.NewTodoServiceClient(conn)
}

func NewTodoServiceConnection(conf *config.Config) (TodoConnection, func(), error) {
	return grpcutil.Dial(conf.TodoAddr, true)
}
