//go:build wireinject

package registry

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	todo_v1 "github.com/phamquanandpad/training-project/grpc/go/todo/todo/v1"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/config"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/interceptor/datavalidator"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/infrastructure/datastore"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
)

func NewGRPCServer(
	todoServiceServer todo_v1.TodoServiceServer,
	dataValidator usecase.DataValidator,
	conf *config.Config,
) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				datavalidator.UnaryServerInterceptor(dataValidator),
			),
		),

		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:              150 * time.Second,
			MaxConnectionIdle: 10 * time.Minute,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime: 1 * time.Minute,
		}),
	)

	todo_v1.RegisterTodoServiceServer(server, todoServiceServer)

	return server
}

func newDBConfig(conf *config.Config) *config.DBConfig {
	return &config.DBConfig{
		DBHost: conf.DBHost,
		DBPort: conf.DBPort,
		DBUser: conf.DBUser,
		DBPass: conf.DBPass,
		DBName: conf.DBName,
	}
}

func InitializeServer(
	conf *config.Config,
	validate *validator.Validate,
) (*grpc.Server, func(), error) {
	wire.Build(
		NewGRPCServer,

		requestbinder.NewRequestBinder,

		handler.NewTodoService,

		service.NewTodoHelper,
		service.NewTodoUpdater,
		service.NewTodoCreator,
		service.NewTodoGetter,
		service.NewTodoLister,
		service.NewTodoDeleter,
		service.NewUserGetter,
		service.NewUserCreator,
		service.NewDataValidator,

		datastore.NewTodoReader,
		datastore.NewTodoWriter,
		datastore.NewUserReader,
		datastore.NewUserWriter,

		newDBConfig,
		datastore.NewConnectionBinder,
		datastore.NewTodoSQLHandler,
	)

	return nil, nil, nil
}
