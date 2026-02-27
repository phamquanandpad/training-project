package datastore

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
)

type userReader struct{}

func NewUserReader() gateway.UserQueriesGateway {
	return &userReader{}
}

func (r *userReader) GetUser(ctx context.Context, userID todo.UserID) (*todo.User, error) {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return nil, err
	}
	db := tx.WithContext(ctx)

	var user todo.User
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
