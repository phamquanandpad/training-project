//go:build wireinject

package registry

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	todo_service "github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/todo"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/config"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/handler/requestbinder"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/datastore"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/jwt"
)

func NewGRPCServer(
	authServiceServer auth_v1.AuthServiceServer,
	conf *config.Config,
) *grpc.Server {

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:              150 * time.Second,
			MaxConnectionIdle: 10 * time.Minute,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime: 1 * time.Minute,
		}),
	)

	auth_v1.RegisterAuthServiceServer(server, authServiceServer)

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

func newJwtConfig(conf *config.Config) *config.JwtConfig {
	return &config.JwtConfig{
		AccessTokenSecret:         conf.AccessTokenSecret,
		RefreshTokenSecret:        conf.RefreshTokenSecret,
		AccessTokenExpiresSecond:  conf.AccessTokenExpiresSecond,
		RefreshTokenExpiresSecond: conf.RefreshTokenExpiresSecond,
	}
}

func InitializeServer(
	c *config.Config,
	v *validator.Validate,
) (*grpc.Server, func(), error) {
	wire.Build(
		NewGRPCServer,

		requestbinder.NewRequestBinder,

		handler.NewAuthService,

		service.NewUserLogin,
		service.NewUserRegister,
		service.NewTokenVerify,
		service.NewTokenRefresh,

		datastore.NewUserReader,
		datastore.NewUserWriter,

		newJwtConfig,
		jwt.NewTokenGenerator,
		jwt.NewTokenVerifier,

		todo_service.NewUserWriter,

		newDBConfig,
		datastore.NewConnectionBinder,
		datastore.NewTodoSQLHandler,

		todo_service.NewTodoServiceClient,
		todo_service.NewTodoServiceConnection,
	)

	return nil, nil, nil
}
