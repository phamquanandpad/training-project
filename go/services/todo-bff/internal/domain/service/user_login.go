package service

import (
	"context"
	"fmt"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase/output"
)

type userLogin struct {
	authReaderGateway gateway.AuthReaderGateway
}

func NewUserLogin(
	authReaderGateway gateway.AuthReaderGateway,
) usecase.UserLogin {
	return &userLogin{
		authReaderGateway: authReaderGateway,
	}
}

func (s userLogin) Login(
	ctx context.Context,
	in *input.UserLogin,
) (*output.UserLogin, error) {
	tokens, err := s.authReaderGateway.Login(
		ctx,
		in.Email,
		in.Password,
	)
	if tokens == nil {
		return nil, app_errors.NewNotFoundError(
			"userLogin.Login",
			fmt.Errorf("user not found"),
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}
	if err != nil {
		return nil, app_errors.NewInternalError(
			"userLogin.Login",
			err,
		)
	}

	return &output.UserLogin{
		Tokens: tokens,
	}, nil
}
