package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type dataValidator struct {
	dbConnBinder       gateway.Binder
	userQueriesGateway gateway.UserQueriesGateway
}

func NewDataValidator(
	dbConnBinder gateway.Binder,
	userQueriesGateway gateway.UserQueriesGateway,
) usecase.DataValidator {
	return &dataValidator{
		dbConnBinder:       dbConnBinder,
		userQueriesGateway: userQueriesGateway,
	}
}

func (s dataValidator) ValidateUserRequest(
	ctx context.Context,
	in *input.UserRequestValidator,
) (*output.UserRequestValidator, error) {
	ctx = s.dbConnBinder.Bind(ctx)

	user, err := s.userQueriesGateway.GetUser(ctx, in.UserID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"dataValidator.ValidateUserRequest",
			err,
		)
	}

	if user == nil {
		return nil, app_errors.NewNotFoundError(
			"dataValidator.ValidateUserRequest", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}

	return (*output.UserRequestValidator)(user), nil
}
