package middleware

import (
	"context"
	"net/http"
)

type Loaders struct{}

func (m middleware) WithGraphql() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return graphMiddleware(h)
	}

}

func graphMiddleware(next http.Handler) http.Handler {
	loadersKey := contextKey("loaders")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{})
		ctx = context.WithValue(ctx, responseWriterKey, w)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
