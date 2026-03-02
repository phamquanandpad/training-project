package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/auth/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

type tokenRefresh struct {
	dbConnBinder       gateway.Binder
	jwtGenerateGateway gateway.JwtGenerateGateway
	jwtVerifyGateway   gateway.JwtVerifyGateway
	userQueriesGateway gateway.UserQueriesGateway
}

func NewTokenRefresh(
	dbConnBinder gateway.Binder,
	jwtGenerateGateway gateway.JwtGenerateGateway,
	jwtVerifyGateway gateway.JwtVerifyGateway,
	userQueriesGateway gateway.UserQueriesGateway,
) usecase.TokenRefresh {
	return &tokenRefresh{
		dbConnBinder:       dbConnBinder,
		jwtGenerateGateway: jwtGenerateGateway,
		jwtVerifyGateway:   jwtVerifyGateway,
		userQueriesGateway: userQueriesGateway,
	}
}

func (s tokenRefresh) RefreshToken(
	ctx context.Context,
	in *input.TokenRefresh,
) (*output.TokenRefresh, error) {
	ctx = s.dbConnBinder.Bind(ctx)

	userID, err := s.jwtVerifyGateway.VerifyRefreshToken(in.RefreshToken)
	if err != nil {
		return nil, app_errors.NewAuthNError(
			"tokenRefresh.Refresh",
			err,
			&app_errors.LocalizedMessage{JaMessage: app_errors.AuthNJaMessage},
		)
	}

	user, err := s.userQueriesGateway.GetUserByID(ctx, userID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"tokenRefresh.Refresh",
			err,
		)
	}

	if user == nil {
		return nil, app_errors.NewNotFoundError(
			"tokenRefresh.Refresh", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}

	access_token, access_token_expire_second, err := s.jwtGenerateGateway.GenerateAccessToken(userID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"tokenRefresh.Refresh",
			err,
		)
	}

	return &output.TokenRefresh{
		AccessToken:              access_token,
		AccessTokenExpiresSecond: access_token_expire_second,
	}, nil
}
