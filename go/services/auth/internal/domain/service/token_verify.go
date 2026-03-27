package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/auth/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

type tokenVerify struct {
	dbConnBinder       gateway.Binder
	jwtVerifyGateway   gateway.JwtVerifyGateway
	userQueriesGateway gateway.UserQueriesGateway
}

func NewTokenVerify(
	dbConnBinder gateway.Binder,
	jwtVerifyGateway gateway.JwtVerifyGateway,
	userQueriesGateway gateway.UserQueriesGateway,
) usecase.TokenVerify {
	return &tokenVerify{
		dbConnBinder:       dbConnBinder,
		jwtVerifyGateway:   jwtVerifyGateway,
		userQueriesGateway: userQueriesGateway,
	}
}

func (s tokenVerify) VerifyToken(
	ctx context.Context,
	in *input.TokenVerify,
) (*output.TokenVerify, error) {
	ctx = s.dbConnBinder.Bind(ctx)

	userID, err := s.jwtVerifyGateway.VerifyAccessToken(in.AccessToken)
	if err != nil {
		return nil, app_errors.NewAuthNError(
			"tokenVerify.VerifyToken",
			err,
			&app_errors.LocalizedMessage{JaMessage: app_errors.AuthNJaMessage},
		)
	}

	user, err := s.userQueriesGateway.GetUserByID(ctx, userID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"tokenVerify.VerifyToken",
			err,
		)
	}

	if user == nil {
		return nil, app_errors.NewNotFoundError(
			"tokenVerify.VerifyToken", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}

	return &output.TokenVerify{
		UserID: userID,
	}, nil
}
