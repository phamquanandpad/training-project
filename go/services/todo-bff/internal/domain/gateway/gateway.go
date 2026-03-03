package gateway

import (
	"context"

	auth_models "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/todo"
)

type Binder interface {
	Bind(ctx context.Context) context.Context
}

type TodoQueriesGateway interface {
	Get(
		ctx context.Context,
		userAttributes todo.UserAttributes,
		todoID todo.TodoID,
	) (*todo.Todo, error)

	List(
		ctx context.Context,
		userAttributes todo.UserAttributes,
		limit int,
		offset int,
	) (todos []*todo.Todo, total int, err error)
}

type TodoCommandsGateway interface {
	Create(
		ctx context.Context,
		userAttributes todo.UserAttributes,
		newTodo todo.NewTodo,
	) (*todo.Todo, error)

	Update(
		ctx context.Context,
		userAttributes todo.UserAttributes,
		todoID todo.TodoID,
		updateTodo todo.UpdateTodo,
	) (*todo.Todo, error)

	Delete(
		ctx context.Context,
		userAttributes todo.UserAttributes,
		todoID todo.TodoID,
	) error
}

type AuthReaderGateway interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (*auth_models.Tokens, error)

	Register(
		ctx context.Context,
		username string,
		email string,
		password string,
	) error

	RefreshToken(
		ctx context.Context,
		refreshToken string,
	) (*auth_models.AccessToken, error)

	VerifyToken(
		ctx context.Context,
		accessToken string,
	) (*auth_models.UserID, error)
}

type SessionQueriesGateway interface {
	GetAccessToken(
		ctx context.Context,
		sessionID string,
	) (*auth_models.AccessToken, error)
}

type SessionCommandsGateway interface {
	SetAccessToken(
		ctx context.Context,
		sessionID string,
		refreshToken string,
		expireSecond int64,
	) error

	DeleteSession(
		ctx context.Context,
		sessionID string,
	) error
}
