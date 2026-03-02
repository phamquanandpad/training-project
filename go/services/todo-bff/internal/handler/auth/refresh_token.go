package auth

import (
	"fmt"
	"net/http"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/middleware"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
)

func (a *authentication) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := middleware.RefreshTokenFromContext(r.Context())
	if refreshToken == "" {
		renderError(w, r, http.StatusUnauthorized, app_errors.NewAuthNError(
			"auth.RefreshToken",
			fmt.Errorf("refresh token not found"),
			&app_errors.LocalizedMessage{JaMessage: app_errors.AuthNJaMessage},
		))
		return
	}

	out, err := a.tokenRefresh.RefreshToken(r.Context(), &input.TokenRefresh{
		RefreshToken: refreshToken,
	})
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.AccessTokenCookieKey,
		Value:    out.AccessToken.Token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}
