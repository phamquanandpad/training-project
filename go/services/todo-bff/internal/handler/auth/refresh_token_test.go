package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"

	mock_usecase "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/mock"

	auth_model "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler/auth"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/middleware"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type PrepareRefreshTokenFields struct {
	mockTokenRefresh *mock_usecase.MockTokenRefresh
}

type RefreshTokenTestcase struct {
	prepare        func(f *PrepareRefreshTokenFields)
	refreshToken   string
	expectedStatus int
	checkCookies   bool
}

func Test_authentication_RefreshToken(t *testing.T) {
	t.Parallel()

	testTables := map[string]RefreshTokenTestcase{
		"Refresh token successfully": {
			prepare: func(f *PrepareRefreshTokenFields) {
				f.mockTokenRefresh.
					EXPECT().
					RefreshToken(gomock.Any(), &input.TokenRefresh{
						RefreshToken: "valid-refresh-token",
					}).
					Return(&output.TokenRefresh{
						AccessToken: &auth_model.AccessToken{
							Token:   "new-access-token",
							Expires: 3600,
						},
					}, nil).
					Times(1)
			},
			refreshToken:   "valid-refresh-token",
			expectedStatus: http.StatusOK,
			checkCookies:   true,
		},
		"Fail to refresh token when refresh token is missing": {
			prepare:        func(f *PrepareRefreshTokenFields) {},
			refreshToken:   "",
			expectedStatus: http.StatusUnauthorized,
			checkCookies:   false,
		},
		"Fail to refresh token when usecase returns error": {
			prepare: func(f *PrepareRefreshTokenFields) {
				f.mockTokenRefresh.
					EXPECT().
					RefreshToken(gomock.Any(), &input.TokenRefresh{
						RefreshToken: "expired-refresh-token",
					}).
					Return(nil, errors.New("token expired")).
					Times(1)
			},
			refreshToken:   "expired-refresh-token",
			expectedStatus: http.StatusInternalServerError,
			checkCookies:   false,
		},
	}

	for name, tt := range testTables {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockTokenRefresh := mock_usecase.NewMockTokenRefresh(ctrl)

			tt.prepare(&PrepareRefreshTokenFields{
				mockTokenRefresh: mockTokenRefresh,
			})

			authenticator := auth.NewAuthentication(
				context.Background(),
				mockTokenRefresh,
			)

			req := httptest.NewRequest(http.MethodPost, "/refresh-token", nil)

			if tt.refreshToken != "" {
				cookies := &middleware.Cookies{
					RefreshToken: middleware.CookiesContextKey(tt.refreshToken),
				}
				ctx := context.WithValue(req.Context(), middleware.CookieKey, cookies)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()

			authenticator.RefreshToken(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("RefreshToken() status = %d, want %d", res.StatusCode, tt.expectedStatus)
			}

			if tt.checkCookies {
				cookieMap := make(map[string]string)
				for _, c := range res.Cookies() {
					cookieMap[c.Name] = c.Value
				}
				if cookieMap[middleware.AccessTokenCookieKey] != "new-access-token" {
					t.Errorf("RefreshToken() access_token cookie = %q, want %q", cookieMap[middleware.AccessTokenCookieKey], "new-access-token")
				}
			}
		})
	}
}
