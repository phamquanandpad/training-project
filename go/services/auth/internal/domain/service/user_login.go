package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	app_errors "github.com/phamquanandpad/training-project/go/services/auth/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

type userLogin struct {
	dbConnBinder       gateway.Binder
	userQueriesGateway gateway.UserQueriesGateway
	jwtGenerateGateway gateway.JwtGenerateGateway
}

func NewUserLogin(
	dbConnBinder gateway.Binder,
	userQueriesGateway gateway.UserQueriesGateway,
	jwtGenerateGateway gateway.JwtGenerateGateway,
) usecase.UserLogin {
	return &userLogin{
		dbConnBinder:       dbConnBinder,
		userQueriesGateway: userQueriesGateway,
		jwtGenerateGateway: jwtGenerateGateway,
	}
}

func (s userLogin) Login(
	ctx context.Context,
	in *input.UserLogin,
) (*output.UserLogin, error) {
	ctx = s.dbConnBinder.Bind(ctx)

	user, err := s.userQueriesGateway.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"userLogin.Login",
			err,
		)
	}

	if user == nil {
		return nil, app_errors.NewNotFoundError(
			"userLogin.Login", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(in.Password),
	); err != nil {
		return nil, app_errors.NewAuthNError(
			"userLogin.Login",
			err,
			&app_errors.LocalizedMessage{JaMessage: app_errors.AuthNJaMessage},
		)
	}

	access_token, access_token_expire_second, err := s.jwtGenerateGateway.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"userLogin.Login",
			err,
		)
	}

	refresh_token, refresh_token_expire_second, err := s.jwtGenerateGateway.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"userLogin.Login",
			err,
		)
	}

	return &output.UserLogin{
		UserID:                    user.ID,
		AccessToken:               access_token,
		AccessTokenExpiresSecond:  access_token_expire_second,
		RefreshToken:              refresh_token,
		RefreshTokenExpiresSecond: refresh_token_expire_second,
	}, nil
}
