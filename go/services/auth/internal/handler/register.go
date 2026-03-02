package handler

import (
	"context"

	auth_v1 "github.com/phamquanandpad/training-project/grpc/go/auth/v1"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
)

func (h *authService) Register(
	ctx context.Context,
	req *auth_v1.RegisterRequest,
) (*auth_v1.RegisterResponse, error) {
	var in input.UserRegister

	if err := h.requestBinder.Bind(ctx, req, &in); err != nil {
		return nil, err
	}

	if err := h.userRegister.Register(ctx, &in); err != nil {
		return nil, err
	}

	return &auth_v1.RegisterResponse{}, nil
}
