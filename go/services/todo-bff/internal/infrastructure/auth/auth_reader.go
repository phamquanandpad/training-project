package auth

import (
	"context"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	auth_models "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/model/auth"
	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
)

type authReader struct {
	authServiceClient auth_v1.AuthServiceClient
}

func NewAuthReader(
	client auth_v1.AuthServiceClient,
) gateway.AuthReaderGateway {
	return &authReader{
		authServiceClient: client,
	}
}

func (r *authReader) Login(
	ctx context.Context,
	email string,
	password string,
) (*auth_models.Tokens, error) {
	req := &auth_v1.LoginRequest{
		Email:    email,
		Password: password,
	}

	res, err := r.authServiceClient.Login(ctx, req)
	if err != nil {
		return nil, app_errors.GrpcStatusToAppError(err)
	}

	return &auth_models.Tokens{
		AccessToken: auth_models.AccessToken{
			Token:   res.AccessToken,
			Expires: res.AccessTokenExpiresSecond,
		},
		RefreshToken: auth_models.RefreshToken{
			Token:   res.RefreshToken,
			Expires: res.RefreshTokenExpiresSecond,
		},
	}, nil
}

func (r *authReader) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
) error {
	req := &auth_v1.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	_, err := r.authServiceClient.Register(ctx, req)
	if err != nil {
		return app_errors.GrpcStatusToAppError(err)
	}
	return nil
}

func (r *authReader) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*auth_models.AccessToken, error) {
	req := &auth_v1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	res, err := r.authServiceClient.RefreshToken(ctx, req)
	if err != nil {
		return nil, app_errors.GrpcStatusToAppError(err)
	}

	return &auth_models.AccessToken{
		Token:   res.AccessToken,
		Expires: res.AccessTokenExpiresSecond,
	}, nil
}

func (r *authReader) VerifyToken(
	ctx context.Context,
	accessToken string,
) (*auth_models.UserID, error) {
	req := &auth_v1.VerifyTokenRequest{
		AccessToken: accessToken,
	}

	res, err := r.authServiceClient.VerifyToken(ctx, req)
	if err != nil {
		return nil, app_errors.GrpcStatusToAppError(err)
	}

	return auth_models.NewUserID(res.UserId), nil
}
