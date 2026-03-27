package middleware

import (
	"net/http"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/config"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
)

type Middleware interface {
	WithCors() func(http.Handler) http.Handler
	WithCookies() func(http.Handler) http.Handler
	WithAuth(tokenVerify usecase.TokenVerify) func(http.Handler) http.Handler
	WithGraphql() func(http.Handler) http.Handler
}

type middleware struct {
	c *config.Config
}

func NewMiddleware(c *config.Config) Middleware {
	return middleware{
		c: c,
	}
}
