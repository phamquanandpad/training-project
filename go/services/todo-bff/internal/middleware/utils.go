package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	auth_models "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/auth"
)

func renderError(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
	// nolint: errcheck
	w.Write([]byte(msg))
}

func ResponseWriterFromContext(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	return w, ok
}

// UserIDFromContext lấy UserID đã được inject bởi WithAuth middleware.
func UserIDFromContext(ctx context.Context) *auth_models.UserID {
	userID, _ := ctx.Value(UserIDKey).(*auth_models.UserID)
	return userID
}

// Helper function to check if a word is a GraphQL keyword
func isGraphQLKeyword(word string) bool {
	keywords := map[string]bool{
		"query":        true,
		"mutation":     true,
		"subscription": true,
		"fragment":     true,
		"on":           true,
		"true":         true,
		"false":        true,
		"null":         true,
		"if":           true,
		"else":         true,
		"type":         true,
		"interface":    true,
		"union":        true,
		"enum":         true,
		"input":        true,
		"extend":       true,
		"schema":       true,
		"directive":    true,
		"scalar":       true,
		"implements":   true,
	}
	return keywords[strings.ToLower(word)]
}
