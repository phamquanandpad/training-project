package handler

import (
	"context"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/output"
)

func (h *authService) VerifyToken(
	ctx context.Context,
	req *auth_v1.VerifyTokenRequest,
) (*auth_v1.VerifyTokenResponse, error) {
	var in input.TokenVerify

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	out, err := h.tokenVerify.VerifyToken(ctx, &in)
	if err != nil {
		return nil, err
	}

	return toVerifyTokenResponse(out), nil
}

func toVerifyTokenResponse(out *output.TokenVerify) *auth_v1.VerifyTokenResponse {
	return &auth_v1.VerifyTokenResponse{
		UserId: out.UserID.Int64(),
	}
}
