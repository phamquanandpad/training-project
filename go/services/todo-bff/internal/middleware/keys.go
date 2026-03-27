package middleware

type contextKey string

const (
	responseWriterKey contextKey = "responseWriter"
	UserIDKey         contextKey = "userID"

	AccessTokenCookieKey  = "access_token"
	RefreshTokenCookieKey = "refresh_token"
)
