package auth

import (
	"context"
	"net/http"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
)

type Authenticator interface {
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type authentication struct {
	tokenRefresh usecase.TokenRefresh
}

func NewAuthentication(
	ctx context.Context,
	tokenRefresh usecase.TokenRefresh,
) Authenticator {
	return &authentication{
		tokenRefresh: tokenRefresh,
	}
}
