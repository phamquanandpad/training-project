package middleware

import (
	"context"
	"net/http"
)

type CookiesContextKey string

type Cookies struct {
	AccessToken  CookiesContextKey
	RefreshToken CookiesContextKey
}

const CookieKey CookiesContextKey = "cookies"

func (m middleware) WithCookies() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var accessTokenValue string
			var refreshTokenValue string

			if accessToken, err := r.Cookie(AccessTokenCookieKey); err == nil && accessToken != nil {
				accessTokenValue = accessToken.Value
			}

			if refreshToken, err := r.Cookie(RefreshTokenCookieKey); err == nil && refreshToken != nil {
				refreshTokenValue = refreshToken.Value
			}

			cookies := &Cookies{
				AccessToken:  CookiesContextKey(accessTokenValue),
				RefreshToken: CookiesContextKey(refreshTokenValue),
			}

			ctx := context.WithValue(r.Context(), CookieKey, cookies)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AccessTokenFromContext(ctx context.Context) string {
	cookies, ok := ctx.Value(CookieKey).(*Cookies)
	if !ok || cookies == nil {
		return ""
	}
	return string(cookies.AccessToken)
}

func RefreshTokenFromContext(ctx context.Context) string {
	cookies, ok := ctx.Value(CookieKey).(*Cookies)
	if !ok || cookies == nil {
		return ""
	}
	return string(cookies.RefreshToken)
}
