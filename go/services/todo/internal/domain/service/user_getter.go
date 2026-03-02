package service

import (
	"context"

	app_errors "github.com/phamquanandpad/training-project/go/services/todo/internal/errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/input"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/usecase/output"
)

type userGetter struct {
	dbConnBinder       gateway.Binder
	userQueriesGateway gateway.UserQueriesGateway
}

func NewUserGetter(
	dbConnBinder gateway.Binder,
	userQueriesGateway gateway.UserQueriesGateway,
) usecase.UserGetter {
	return &userGetter{
		dbConnBinder:       dbConnBinder,
		userQueriesGateway: userQueriesGateway,
	}
}

func (s userGetter) Get(
	ctx context.Context,
	in *input.UserGetter,
) (*output.UserGetter, error) {
	ctx = s.dbConnBinder.Bind(ctx)

	user, err := s.userQueriesGateway.GetUser(ctx, in.UserID)
	if err != nil {
		return nil, app_errors.NewInternalError(
			"userGetter.GetUser",
			err,
		)
	}

	if user == nil {
		return nil, app_errors.NewNotFoundError(
			"userGetter.GetUser", nil,
			&app_errors.LocalizedMessage{JaMessage: app_errors.NotFoundJaMessage},
		)
	}

	return &output.UserGetter{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
