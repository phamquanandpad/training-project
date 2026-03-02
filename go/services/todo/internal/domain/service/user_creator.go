package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type userCreator struct {
	dbConnBinder        gateway.Binder
	userCommandsGateway gateway.UserCommandsGateway
}

func NewUserCreator(
	dbConnBinder gateway.Binder,
	userCommandsGateway gateway.UserCommandsGateway,
) usecase.UserCreator {
	return &userCreator{
		dbConnBinder:        dbConnBinder,
		userCommandsGateway: userCommandsGateway,
	}
}

func (s userCreator) Create(
	ctx context.Context,
	in *input.UserCreator,
) (*output.UserCreator, error) {
	ctx = s.dbConnBinder.Bind(ctx)

	newUser := todo.NewUser{
		ID:        in.User.ID,
		Username:  in.User.Username,
		Email:     in.User.Email,
		CreatedAt: in.User.CreatedAt,
		UpdatedAt: in.User.UpdatedAt,
	}
	user, err := s.userCommandsGateway.CreateUser(ctx, newUser)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"userCreator.CreateUser",
			err,
		)
	}

	return &output.UserCreator{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
