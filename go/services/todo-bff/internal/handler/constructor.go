package handler

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/config"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler/auth"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler/graph"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler/graph/generated"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/middleware"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
)

func NewHTTPServer(
	conf *config.Config,
	v *validator.Validate,
	middle middleware.Middleware,
	todoGetter usecase.TodoGetter,
	todoLister usecase.TodoLister,
	todoCreator usecase.TodoCreator,
	todoUpdater usecase.TodoUpdater,
	todoDeleter usecase.TodoDeleter,
	userLogin usecase.UserLogin,
	userRegister usecase.UserRegister,
	tokenRefresh usecase.TokenRefresh,
	tokenVerify usecase.TokenVerify,
	authenticator auth.Authenticator,
) http.Handler {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		router.Route("/graphql", func(r chi.Router) {
			r.Use(
				middle.WithCors(),
				middle.WithAuth(tokenVerify),
				middle.WithGraphql(),
			)

			graphConf := graph.New(
				todoGetter,
				todoLister,
				todoCreator,
				todoUpdater,
				todoDeleter,
			)

			srv := handler.New(generated.NewExecutableSchema(graphConf))

			srv.AddTransport(transport.POST{})
			srv.AddTransport(transport.Options{})
			srv.AddTransport(transport.GET{})

			srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
			srv.SetErrorPresenter(middleware.ErrorPresenter())
			srv.SetRecoverFunc(middleware.RecoverPanicError())

			srv.Use(extension.Introspection{})
			srv.Use(extension.AutomaticPersistedQuery{
				Cache: lru.New[string](100),
			})

			r.Handle("/", srv)
		})

		router.Route("/graphql/auth", func(r chi.Router) {
			r.Use(
				middle.WithCors(),
				middle.WithCookies(),
				middle.WithGraphql(),
			)

			graphConf := graph.NewAuth(
				userLogin,
				userRegister,
			)

			srv := handler.New(generated.NewExecutableSchema(graphConf))

			srv.AddTransport(transport.POST{})
			srv.AddTransport(transport.Options{})
			srv.AddTransport(transport.GET{})

			srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
			srv.SetErrorPresenter(middleware.ErrorPresenter())
			srv.SetRecoverFunc(middleware.RecoverPanicError())

			srv.Use(extension.Introspection{})
			srv.Use(extension.AutomaticPersistedQuery{
				Cache: lru.New[string](100),
			})

			r.Handle("/", srv)
		})

		router.Route("/auth", func(r chi.Router) {
			r.Use(
				middle.WithCors(),
				middle.WithCookies(),
			)
			r.Post("/refresh", authenticator.RefreshToken)
		})
	})

	if conf.Env != config.EnvProduction {
		router.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	}
	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "OK\n")
	})

	return router
}
