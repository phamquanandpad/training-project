package datastore

import (
	"context"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
)

type userWriter struct{}

func NewUserWriter() gateway.UserCommandsGateway {
	return &userWriter{}
}

func (w *userWriter) CreateUser(ctx context.Context, newUser auth_models.NewUser) (*auth_models.User, error) {
	tx, err := ExtractAuthDB(ctx)
	if err != nil {
		return nil, err
	}
	db := tx.WithContext(ctx)

	createdUser := auth_models.User{
		Email:    newUser.Email,
		Password: newUser.Password,
	}

	err = db.
		Create(&createdUser).
		Error
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}
