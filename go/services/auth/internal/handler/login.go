package handler

import (
	"context"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

func (h *authService) Login(
	ctx context.Context,
	req *auth_v1.LoginRequest,
) (*auth_v1.LoginResponse, error) {
	var in input.UserLogin

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	out, err := h.userLogin.Login(ctx, &in)
	if err != nil {
		return nil, err
	}

	return toLoginResponse(out), nil
}

func toLoginResponse(out *output.UserLogin) *auth_v1.LoginResponse {
	return &auth_v1.LoginResponse{
		AccessToken:               out.AccessToken,
		AccessTokenExpiresSecond:  int64(out.AccessTokenExpiresSecond),
		RefreshToken:              out.RefreshToken,
		RefreshTokenExpiresSecond: out.RefreshTokenExpiresSecond,
	}
}
