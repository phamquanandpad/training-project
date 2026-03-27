package middleware

import (
	"context"
	"errors"
	"net/http"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
)

func (m middleware) WithAuth(tokenVerify usecase.TokenVerify) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(AccessTokenCookieKey)
			if err != nil {
				renderError(w, r, http.StatusUnauthorized, app_errors.NewAuthNError(
					"middleware.WithAuth",
					errors.New("access token cookie not found"),
					&app_errors.LocalizedMessage{JaMessage: app_errors.AuthNJaMessage},
				))
				return
			}

			// 2. Verify access token qua usecase
			out, err := tokenVerify.VerifyToken(r.Context(), &input.TokenVerify{
				AccessToken: cookie.Value,
			})
			if err != nil {
				renderError(w, r, http.StatusUnauthorized, app_errors.NewAuthNError(
					"middleware.WithAuth",
					err,
					&app_errors.LocalizedMessage{JaMessage: app_errors.AuthNJaMessage},
				))
				return
			}

			// 3. Đưa UserID vào context cho downstream handlers
			ctx := context.WithValue(r.Context(), UserIDKey, out.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
