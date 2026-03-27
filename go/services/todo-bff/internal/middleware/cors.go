package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func (m middleware) WithCors() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   m.c.AllowOrigins,
		AllowedHeaders:   m.c.AllowHeaders,
		AllowCredentials: true,
		Debug:            m.c.Env == "local",
	}).Handler
}
