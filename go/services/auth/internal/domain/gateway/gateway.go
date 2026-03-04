package gateway

import (
	"context"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/todo"
)

type Binder interface {
	Bind(context.Context) context.Context
}

type UserQueriesGateway interface {
	GetUserByID(ctx context.Context, userID auth_models.UserID) (*auth_models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*auth_models.User, error)
}

type UserCommandsGateway interface {
	CreateUser(ctx context.Context, newUser auth_models.NewUser) (*auth_models.User, error)
}

type JwtGenerateGateway interface {
	GenerateAccessToken(userID auth_models.UserID) (string, int64, error)
	GenerateRefreshToken(userID auth_models.UserID) (string, int64, error)
}

type JwtVerifyGateway interface {
	VerifyAccessToken(token string) (auth_models.UserID, error)
	VerifyRefreshToken(token string) (auth_models.UserID, error)
}

type UserServiceCommandsGateway interface {
	CreateUser(ctx context.Context, newUser todo.NewUser) (*todo.User, error)
}
