package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type tokenVerify struct {
	authReaderGateway gateway.AuthReaderGateway
}

func NewTokenVerify(
	authReaderGateway gateway.AuthReaderGateway,
) usecase.TokenVerify {
	return &tokenVerify{
		authReaderGateway: authReaderGateway,
	}
}

func (s tokenVerify) VerifyToken(
	ctx context.Context,
	in *input.TokenVerify,
) (*output.TokenVerify, error) {
	userID, err := s.authReaderGateway.VerifyToken(
		ctx,
		in.AccessToken,
	)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"tokenVerify.VerifyToken",
			err,
		)
	}

	return &output.TokenVerify{
		UserID: userID,
	}, nil
}
