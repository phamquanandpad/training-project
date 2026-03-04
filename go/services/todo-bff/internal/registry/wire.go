//go:build wireinject

package registry

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"

	auth_handler "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler/auth"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/config"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/service"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/infrastructure/auth"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/infrastructure/todo"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/middleware"
)

func InitHttpServer(ctx context.Context, c *config.Config, v *validator.Validate) (http.Handler, func(), error) {
	wire.Build(
		handler.NewHTTPServer,
		middleware.NewMiddleware,
		auth_handler.NewAuthentication,

		service.NewUserLogin,
		service.NewUserRegister,
		service.NewTokenRefresh,
		service.NewTokenVerify,
		service.NewTodoGetter,
		service.NewTodoLister,
		service.NewTodoCreator,
		service.NewTodoUpdater,
		service.NewTodoDeleter,

		auth.NewAuthReader,

		todo.NewTodoReader,
		todo.NewTodoWriter,

		auth.NewAuthServiceConnection,
		auth.NewAuthServiceClient,

		todo.NewTodoServiceConnection,
		todo.NewTodoServiceClient,
	)
	return nil, nil, nil
}
