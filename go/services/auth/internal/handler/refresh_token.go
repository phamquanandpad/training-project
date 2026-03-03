package handler

import (
	"context"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

func (h *authService) RefreshToken(
	ctx context.Context,
	req *auth_v1.RefreshTokenRequest,
) (*auth_v1.RefreshTokenResponse, error) {
	var in input.TokenRefresh

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	out, err := h.tokenRefresh.RefreshToken(ctx, &in)
	if err != nil {
		return nil, err
	}

	return toRefreshTokenResponse(out), nil
}

func toRefreshTokenResponse(out *output.TokenRefresh) *auth_v1.RefreshTokenResponse {
	return &auth_v1.RefreshTokenResponse{
		AccessToken:               out.AccessToken,
		AccessTokenExpireDuration: out.AccessTokenExpireDuration,
	}
}
