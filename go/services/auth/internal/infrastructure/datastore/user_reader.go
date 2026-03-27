package datastore

import (
	"context"
	"errors"

	"gorm.io/gorm"

	auth_models "github.com/phamquanandpad/training-project/go/services/auth/internal/domain/model/auth"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/domain/gateway"
)

type userReader struct{}

func NewUserReader() gateway.UserQueriesGateway {
	return &userReader{}
}

func (r *userReader) GetUserByID(ctx context.Context, userID auth_models.UserID) (*auth_models.User, error) {
	tx, err := ExtractAuthDB(ctx)
	if err != nil {
		return nil, err
	}
	db := tx.WithContext(ctx)

	var user auth_models.User
	err = db.
		Where("id = ?", userID).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userReader) GetUserByEmail(ctx context.Context, email string) (*auth_models.User, error) {
	tx, err := ExtractAuthDB(ctx)
	if err != nil {
		return nil, err
	}
	db := tx.WithContext(ctx)

	var user auth_models.User
	err = db.
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
