package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"
	app_errors "github.com/phamquanandpad/training-project/go/services/auth/internal/errors"

	"github.com/phamquanandpad/training-project/go/pkg/cast"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/todo"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/usecase/input"
)

type userRegister struct {
	dbConnBinder               gateway.Binder
	userCommandsGateway        gateway.UserCommandsGateway
	userQueriesGateway         gateway.UserQueriesGateway
	userServiceCommandsGateway gateway.UserServiceCommandsGateway
}

func NewUserRegister(
	dbConnBinder gateway.Binder,
	userCommandsGateway gateway.UserCommandsGateway,
	userQueriesGateway gateway.UserQueriesGateway,
	userServiceCommandsGateway gateway.UserServiceCommandsGateway,
) usecase.UserRegister {
	return &userRegister{
		dbConnBinder:               dbConnBinder,
		userCommandsGateway:        userCommandsGateway,
		userQueriesGateway:         userQueriesGateway,
		userServiceCommandsGateway: userServiceCommandsGateway,
	}
}

func (s userRegister) Register(
	ctx context.Context,
	in *input.UserRegister,
) error {
	ctx = s.dbConnBinder.Bind(ctx)

	existingUser, err := s.userQueriesGateway.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return app_errors.NewInternalError(
			"userRegister.Register",
			err,
		)
	}
	if existingUser != nil {
		return app_errors.NewAlreadyExistsError(
			"userRegister.Register",
			nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.AlreadyExistsJaMessage},
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(in.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return app_errors.NewInternalError(
			"userRegister.Register",
			err,
		)
	}

	newUser := auth_models.NewUser{
		Email:    in.Email,
		Password: string(hashedPassword),
	}

	user, err := s.userCommandsGateway.CreateUser(ctx, newUser)
	if err != nil {
		return app_errors.NewInternalError(
			"userRegister.Register",
			err,
		)
	}

	_, err = s.userServiceCommandsGateway.CreateUser(ctx, todo.NewUser{
		ID:        *todo.NewUserID(user.ID.Int64()),
		Username:  in.Username,
		Email:     cast.Ptr(in.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		return app_errors.NewInternalError(
			"userRegister.Register",
			err,
		)
	}

	return nil
}
