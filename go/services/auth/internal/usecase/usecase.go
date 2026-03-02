package usecase

import (
	"context"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

type UserLogin interface {
	Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error)
}

type UserRegister interface {
	Register(ctx context.Context, in *input.UserRegister) error
}

type TokenRefresh interface {
	RefreshToken(ctx context.Context, in *input.TokenRefresh) (*output.TokenRefresh, error)
}

type TokenVerify interface {
	VerifyToken(ctx context.Context, in *input.TokenVerify) (*output.TokenVerify, error)
}
