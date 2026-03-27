package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
)

type userRegister struct {
	authReaderGateway gateway.AuthReaderGateway
}

func NewUserRegister(
	authReaderGateway gateway.AuthReaderGateway,
) usecase.UserRegister {
	return &userRegister{
		authReaderGateway: authReaderGateway,
	}
}

func (s userRegister) Register(
	ctx context.Context,
	in *input.UserRegister,
) error {
	err := s.authReaderGateway.Register(
		ctx,
		in.Username,
		in.Email,
		in.Password,
	)
	if err != nil {
		return app_errors.NewInternalError(
			"userRegister.Register",
			err,
		)
	}
	return err
}
