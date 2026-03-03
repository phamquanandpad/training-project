package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type tokenRefresh struct {
	authReaderGateway gateway.AuthReaderGateway
}

func NewTokenRefresh(
	authReaderGateway gateway.AuthReaderGateway,
) usecase.TokenRefresh {
	return &tokenRefresh{
		authReaderGateway: authReaderGateway,
	}
}

func (s tokenRefresh) RefreshToken(
	ctx context.Context,
	in *input.TokenRefresh,
) (*output.TokenRefresh, error) {
	accessToken, err := s.authReaderGateway.RefreshToken(
		ctx,
		in.RefreshToken,
	)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"tokenRefresh.RefreshToken",
			err,
		)
	}

	return &output.TokenRefresh{
		AccessToken: accessToken,
	}, nil
}
